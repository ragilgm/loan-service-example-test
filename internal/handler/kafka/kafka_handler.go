package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/test/loan-service/internal/consts"
	"github.com/test/loan-service/internal/dto/message"
	"github.com/test/loan-service/internal/service"
	"go.uber.org/dig"
)

type handler func(msg kafka.Message) error

type kafkaSvc struct {
	dig.In
	kafkaReader    *kafka.Reader
	loanSvc        service.LoanSvc
	loanFundingSvc service.LoanFundingSvc
	handlers       map[string]handler
}

// NewKafkaHandler membuat handler Kafka baru dan memulai konsumsi pesan
func NewKafkaHandler(kafkaReader *kafka.Reader, loanSvc service.LoanSvc, loanFundingSvc service.LoanFundingSvc) error {
	svc := kafkaSvc{
		kafkaReader:    kafkaReader,
		loanSvc:        loanSvc,
		loanFundingSvc: loanFundingSvc,
		handlers:       make(map[string]handler),
	}

	// register handler
	svc.register(string(consts.ApprovalLoanTopic), svc.ApprovalLoanHandler)
	svc.register(string(consts.LoanDisburseTopic), svc.DisburseLoanHandler)
	svc.register(string(consts.FundingProcessTopic), svc.FundingProcessHandler)

	// start consume
	go svc.startConsuming()

	return nil
}

func (svc *kafkaSvc) register(topic string, h handler) {
	svc.handlers[topic] = h
}

func (svc *kafkaSvc) startConsuming() {
	for {
		// read message
		msg, err := svc.kafkaReader.ReadMessage(context.Background())
		if err != nil {
			logrus.Errorf("Error reading message from Kafka: %v", err)
			continue
		}
		svc.handleMessage(msg)
	}
}

func (svc *kafkaSvc) handleMessage(msg kafka.Message) {

	handler, exists := svc.handlers[msg.Topic]
	if !exists {
		logrus.Warnf("No handler found for topic: %s", msg.Topic)
		return
	}

	// Menjalankan handler untuk pesan yang diterima
	if err := handler(msg); err != nil {
		logrus.Errorf("Error handling message: %v", err)
	}
}

func (svc *kafkaSvc) ApprovalLoanHandler(msg kafka.Message) error {
	var loanApproval message.UpdateLoanMessage
	if err := json.Unmarshal(msg.Value, &loanApproval); err != nil {
		logrus.Errorf("Error unmarshaling message: %v", err)
		return errors.New("failed to unmarshal loan approval")
	}

	ctx := context.Background()

	// update loan
	err := svc.loanSvc.ApprovalLoan(ctx, loanApproval)
	if err != nil {
		logrus.Errorf("Error processing loan approval: %v", err)
		return errors.New("failed to process loan approval")
	}

	// Log sukses
	logrus.Infof("Loan approval processed successfully: %v", loanApproval)
	return nil
}

func (svc *kafkaSvc) DisburseLoanHandler(msg kafka.Message) error {
	var loanApproval message.UpdateLoanMessage
	if err := json.Unmarshal(msg.Value, &loanApproval); err != nil {
		logrus.Errorf("Error unmarshaling message: %v", err)
		return errors.New("failed to unmarshal loan disburse")
	}

	ctx := context.Background()

	// update loan
	err := svc.loanSvc.DisburseLoan(ctx, loanApproval)
	if err != nil {
		logrus.Errorf("Error processing loan disburse: %v", err)
		return errors.New("failed to process loan disburse")
	}

	// Log sukses
	logrus.Infof("Loan disburse processed successfully: %v", loanApproval)
	return nil
}

func (svc *kafkaSvc) FundingProcessHandler(msg kafka.Message) error {
	var fundingProcess message.FundingProcessMessage
	if err := json.Unmarshal(msg.Value, &fundingProcess); err != nil {
		logrus.Errorf("Error unmarshaling message: %v", err)
		return errors.New("failed to unmarshal loan approval")
	}

	ctx := context.Background()

	// update loan
	err := svc.loanFundingSvc.FundingProcess(ctx, fundingProcess)
	if err != nil {
		logrus.Errorf("Error processing loan approval: %v", err)
		return errors.New("failed to process loan approval")
	}

	// Log sukses
	logrus.Infof("Loan approval processed successfully: %v", fundingProcess)
	return nil
}
