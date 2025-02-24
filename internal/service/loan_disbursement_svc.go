package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"github.com/test/loan-service/internal/consts"
	"github.com/test/loan-service/internal/dto"
	message2 "github.com/test/loan-service/internal/dto/message"
	"github.com/test/loan-service/internal/enum"
	repo "github.com/test/loan-service/internal/repository"
	"github.com/test/loan-service/internal/service/models"
	"github.com/test/loan-service/internal/service/validator"
	"github.com/test/loan-service/internal/utils"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
	"go.uber.org/dig"
	"time"
)

type (
	LoanDisbursementSvc interface {
		Create(context.Context, *dto.LoanDisbursementRequestDTO) error
		Update(ctx context.Context, disbursementID int64, disbursementRequest *dto.UpdateLoanDisbursementRequestDTO) error
		GetByID(ctx context.Context, disbursementID int64) (*dto.LoanDisbursementResponseDTO, error)
		GetAllPage(ctx context.Context, request models.LoanDisbursementRequest) ([]dto.LoanDisbursementResponseDTO, int, error)
	}

	LoanDisbursementSvcImpl struct {
		dig.In
		Repo          repo.LoanDisbursementRepo
		LoanRepo      repo.LoanRepo
		LoanDetailSvc LoanDetailSvc
		KafkaWriter   *kafka.Writer
		Validator     validator.LoanDisbursementValidatorImpl
	}
)

func NewLoanDisbursementSvc(impl LoanDisbursementSvcImpl) LoanDisbursementSvc {
	return &impl
}

func (b *LoanDisbursementSvcImpl) Create(ctx context.Context, disbursementRequest *dto.LoanDisbursementRequestDTO) error {
	log.Printf("Create loan disbursement: LoanID=%d, DisburseAmount=%f", disbursementRequest.LoanID, disbursementRequest.DisburseAmount)

	// Validate request
	err := b.Validator.ValidateCreate(disbursementRequest)
	if err != nil {
		log.Printf("Validation failed: %s", err)
		return err
	}

	var disbursement repo.LoanDisbursement
	err = mapstructure.Decode(disbursementRequest, &disbursement)
	if err != nil {
		log.Printf("Error decoding disbursement request: %v", err)
		return err
	}

	disbursement.DisburseCode = utils.GenerateAlphanumericCode(10)
	disbursement.DisbursementStatus = enum.LoanDisbursementPending
	disbursement.AgreementURL = "http://google.com"
	disbursement.CreatedAt = time.Now()
	disbursement.UpdatedAt = time.Now()

	_, err = b.Repo.Create(ctx, &disbursement)
	if err != nil {
		log.Printf("Error creating loan disbursement in repo: %v", err)
		return err
	}

	log.Printf("Loan disbursement created successfully: DisbursementCode=%s", disbursement.DisburseCode)
	return nil
}

func (b *LoanDisbursementSvcImpl) Update(ctx context.Context, disbursementID int64, disbursementRequest *dto.UpdateLoanDisbursementRequestDTO) error {
	log.Printf("Update loan disbursement: DisbursementID=%d, DisbursementStatus=%v", disbursementID, disbursementRequest.DisbursementStatus)

	disbursement, err := b.Repo.GetByID(ctx, disbursementID)
	if err != nil {
		log.Printf("Error retrieving loan disbursement by ID: %v", err)
		return errors.New("99999")
	}

	if disbursement == nil {
		log.Printf("Loan disbursement not found: DisbursementID=%d", disbursementID)
		return errors.New("10001")

	}

	err = b.Validator.ValidateUpdate(disbursementRequest)
	if err != nil {
		log.Printf("Validation failed: %s", err)
		return err
	}

	valid := b.Validator.ValidateTransitionStatus(disbursement.DisbursementStatus, disbursementRequest.DisbursementStatus)
	if !valid {
		log.Printf("Invalid status transition: CurrentStatus=%v, RequestedStatus=%v", disbursement.DisbursementStatus, disbursementRequest.DisbursementStatus)
		return errors.New("10003")
	}

	disbursement.DisbursementStatus = disbursementRequest.DisbursementStatus
	disbursement.StaffID = &disbursementRequest.StaffID
	disbursement.SignedAgreementURL = &disbursementRequest.SignedAgreementURL

	now := time.Now()
	disbursement.DisburseDate = &now
	disbursement.UpdatedAt = time.Now()

	// Start transaction to update loan
	txnCtx := dbtxn.Begin(&ctx)
	defer func() {
		if err = txnCtx.Commit(); err != nil {
			log.WithError(err).Error("Transaction commit failed")
			txnCtx.AppendError(err)
		}
	}()

	err = b.Repo.Update(ctx, disbursement)
	if err != nil {
		log.Printf("Error updating loan disbursement in repo: %v", err)
		txnCtx.AppendError(err)
		return errors.New("99999")
	}

	// publish for processing loan on going
	err = b.publishDisburseLoan(ctx, disbursement.LoanID, enum.Disbursed)
	if err != nil {
		log.Printf("Error updating loan disbursement in repo: %v", err)
		txnCtx.AppendError(err)
		return errors.New("99999")
	}
	log.Printf("Loan disbursement updated successfully: DisbursementID=%d", disbursementID)
	return nil
}

func (b *LoanDisbursementSvcImpl) publishDisburseLoan(ctx context.Context, loanID int64, status enum.LoanStatus) error {
	// publish kafka for update loan
	req := message2.UpdateLoanMessage{
		LoanID:     loanID,
		LoanStatus: status,
	}

	log.Infof("Marshalling loan update message for approval ID: %d", loanID)
	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Errorf("Failed to marshal loan update message for approval ID: %d: %v", loanID, err)
		return errors.New("99999")
	}

	topic := string(consts.LoanDisburseTopic)
	message := kafka.Message{
		Topic: topic,
		Value: jsonData,
	}

	log.Infof("Sending loan update message to Kafka topic: %s", topic)
	err = b.KafkaWriter.WriteMessages(ctx, message)
	if err != nil {
		log.Errorf("Failed to send loan update message to Kafka topic: %s: %v", topic, err)
		return errors.New("99999")
	}

	log.Infof("Successfully published loan update for approval ID: %d", loanID)
	return nil
}

func (b *LoanDisbursementSvcImpl) GetByID(ctx context.Context, disbursementID int64) (*dto.LoanDisbursementResponseDTO, error) {
	log.Printf("Get loan disbursement by ID: DisbursementID=%d", disbursementID)

	disbursement, err := b.Repo.GetByID(ctx, disbursementID)
	if err != nil {
		log.Printf("Error retrieving loan disbursement by ID: %v", err)
		return nil, err
	}

	if disbursement == nil {
		log.Printf("Loan disbursement not found: DisbursementID=%d", disbursementID)
		return nil, errors.New("loan disbursement not found")
	}

	var disbursementResponse dto.LoanDisbursementResponseDTO
	err = mapstructure.Decode(disbursement, &disbursementResponse)
	if err != nil {
		log.Printf("Error decoding disbursement: %v", err)
		return nil, err
	}

	disbursementResponse.CreatedAt = disbursement.CreatedAt
	disbursementResponse.UpdatedAt = disbursement.UpdatedAt
	disbursementResponse.DeletedAt = disbursement.DeletedAt

	log.Printf("Loan disbursement retrieved successfully: DisbursementID=%d", disbursementID)
	return &disbursementResponse, nil
}

func (b *LoanDisbursementSvcImpl) GetAllPage(ctx context.Context, request models.LoanDisbursementRequest) ([]dto.LoanDisbursementResponseDTO, int, error) {
	log.Printf("Get loan disbursements page: Page=%d, Size=%d", request.Page, request.Size)

	offset := (request.Page - 1) * request.Size

	var repoReq repo.LoanDisbursementRequest
	err := mapstructure.Decode(request, &repoReq)
	if err != nil {
		log.Printf("Error decoding loan disbursement request: %v", err)
		return nil, 0, err
	}
	repoReq.Offset = offset

	disbursements, totalRecords, err := b.Repo.GetAllPage(ctx, repoReq)
	if err != nil {
		log.Printf("Error retrieving loan disbursements from repo: %v", err)
		return nil, 0, err
	}

	disbursementDTOs := []dto.LoanDisbursementResponseDTO{}
	for _, disbursement := range disbursements {
		var disbursementResponse dto.LoanDisbursementResponseDTO
		err = mapstructure.Decode(disbursement, &disbursementResponse)
		if err != nil {
			log.Printf("Error decoding disbursement: %v", err)
			return nil, 0, err
		}

		disbursementResponse.CreatedAt = disbursement.CreatedAt
		disbursementResponse.UpdatedAt = disbursement.UpdatedAt
		disbursementResponse.DeletedAt = disbursement.DeletedAt

		disbursementDTOs = append(disbursementDTOs, disbursementResponse)
	}

	log.Printf("Loan disbursements retrieved successfully: TotalRecords=%d", totalRecords)
	return disbursementDTOs, int(totalRecords), nil
}
