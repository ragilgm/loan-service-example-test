package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/test/loan-service/internal/consts"
	"github.com/test/loan-service/internal/dto"
	message2 "github.com/test/loan-service/internal/dto/message"
	"github.com/test/loan-service/internal/enum"
	repo "github.com/test/loan-service/internal/repository"
	"github.com/test/loan-service/internal/utils"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
	"go.uber.org/dig"
	"sync"
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=$PROJ/internal/generated/mock/mock_$GOPACKAGE/$GOFILE

type (
	LoanFundingSvc interface {
		Create(context.Context, *dto.LoanFundingRequestDTO) error
		FundingProcess(context.Context, message2.FundingProcessMessage) error
		GetByID(ctx context.Context, id int64) (*dto.LoanFundingResponseDTO, error)
		GetByLenderID(ctx context.Context, lenderID int64) ([]dto.LoanFundingResponseDTO, error)
	}

	AsyncResultData struct {
		LoanFunding *repo.LoanFunding
		Loan        *repo.Loan
	}

	LoanFundingSvcImpl struct {
		dig.In
		Repo        repo.LoanFundingRepo
		LoanRepo    repo.LoanRepo
		DisburseSvc LoanDisbursementSvc
		KafkaWriter *kafka.Writer
		MailSvc     EmailSvc
	}
)

func NewLoanFundingSvc(impl LoanFundingSvcImpl) LoanFundingSvc {
	return &impl
}

func (s *LoanFundingSvcImpl) Create(ctx context.Context, request *dto.LoanFundingRequestDTO) (err error) {
	// Validate the LoanFundingRequest
	if errMsg := s.validateRequestLoanFunding(request); errMsg != "" {
		logrus.Warnf("Loan funding validation failed: %s", errMsg)
		return errors.New(errMsg)
	}

	var loan *repo.Loan
	loan, err = s.LoanRepo.GetByID(ctx, request.LoanID)
	if err != nil || loan == nil {
		logrus.Errorf("Loan not found for LoanID %d, error: %v", request.LoanID, err)
		return errors.New("99999")
	}

	if loan.LoanStatus != enum.Approved {
		logrus.Warnf("Loan %d status is not approved", loan.ID)
		return errors.New("10003")
	}

	if loan.FundingDeadline.Before(time.Now()) {
		logrus.Warnf("Loan %d funding deadline has passed", loan.ID)
		return errors.New("10003")
	}

	// Create loan funding
	var loanFunding repo.LoanFunding
	err = mapstructure.Decode(request, &loanFunding)
	if err != nil {
		logrus.Errorf("Failed to map request to LoanFunding struct: %v", err)
		return errors.New("99999")
	}

	loanFunding.LoanOrderNumber = utils.GenerateAlphanumericCode(10)
	loanFunding.LenderID = request.LenderID
	loanFunding.Status = enum.LoanFundingPending
	now := time.Now()
	loanFunding.InvestmentDate = now
	loanFunding.CreatedAt = now
	loanFunding.UpdatedAt = now

	// Log creation attempt
	logrus.Infof("Creating loan funding for LoanID %d, LoanOrderNumber %s", loan.ID, loanFunding.LoanOrderNumber)

	// Create initial loan funding
	_, err = s.Repo.Create(ctx, &loanFunding)
	if err != nil {
		logrus.Errorf("Failed to create loan funding for LoanID %d: %v", request.LoanID, err)
		return errors.New("99999")
	}

	// Log success in funding creation
	logrus.Infof("Loan funding created successfully for LoanOrderNumber %s", loanFunding.LoanOrderNumber)

	// Publish for funding process
	err = s.publishFundingProcess(ctx, loanFunding)
	if err != nil {
		logrus.Errorf("Failed to publish funding process for LoanOrderNumber %s: %v", loanFunding.LoanOrderNumber, err)
		return errors.New("99999")
	}

	logrus.Infof("Loan funding process initiated successfully for LoanID %d", request.LoanID)
	return nil
}

func (s *LoanFundingSvcImpl) GetByID(ctx context.Context, id int64) (*dto.LoanFundingResponseDTO, error) {
	// Get loan funding by ID
	loanFunding, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		// Log error on failed query
		logrus.Errorf("Failed to get loan funding by ID %d: %v", id, err)
		return nil, errors.New("99999")
	}

	if loanFunding == nil {
		logrus.Warnf("Loan funding not found for ID %d", id)
		return nil, errors.New("10001")
	}

	var loanFundingRes dto.LoanFundingResponseDTO
	err = mapstructure.Decode(loanFunding, &loanFundingRes)
	if err != nil {
		logrus.Errorf("Failed to map loan funding to response DTO for ID %d: %v", id, err)
		return nil, errors.New("99999")
	}

	loanFundingRes.CreatedAt = loanFunding.CreatedAt
	loanFundingRes.UpdatedAt = loanFunding.UpdatedAt
	loanFundingRes.DeletedAt = loanFunding.DeletedAt

	logrus.Infof("Loan funding found for ID %d", id)
	return &loanFundingRes, nil
}
func (s *LoanFundingSvcImpl) GetByLenderID(ctx context.Context, lenderID int64) ([]dto.LoanFundingResponseDTO, error) {
	// Get loan funding by lender ID
	loanFundings, err := s.Repo.GetByLenderID(ctx, lenderID)
	if err != nil {
		// Log error on failed query
		logrus.Errorf("Failed to get loan fundings for LenderID %d: %v", lenderID, err)
		return nil, fmt.Errorf("failed to get loan fundings for LenderID %d: %v", lenderID, err)
	}

	var loanFundingResponses []dto.LoanFundingResponseDTO
	for _, loanFunding := range loanFundings {
		var loanFundingRes dto.LoanFundingResponseDTO
		err = mapstructure.Decode(loanFunding, &loanFundingRes)
		if err != nil {
			logrus.Errorf("Failed to map loan funding to response DTO: %v", err)
			return nil, errors.New("99999")
		}
		loanFundingRes.CreatedAt = loanFunding.CreatedAt
		loanFundingRes.UpdatedAt = loanFunding.UpdatedAt
		loanFundingRes.DeletedAt = loanFunding.DeletedAt

		loanFundingResponses = append(loanFundingResponses, loanFundingRes)
	}

	logrus.Infof("Found %d loan fundings for LenderID %d", len(loanFundingResponses), lenderID)
	return loanFundingResponses, nil
}

func (s *LoanFundingSvcImpl) FundingProcess(ctx context.Context, message message2.FundingProcessMessage) error {
	var wg sync.WaitGroup
	resultCh := make(chan AsyncResultData, 2) // Channel dengan buffer untuk 2 hasil

	// Menambahkan 2 goroutine yang harus ditunggu
	wg.Add(2)

	// Menjalankan 2 query secara paralel
	go func() {
		defer wg.Done() // Menandakan bahwa goroutine ini selesai
		res, err := s.LoanRepo.GetByID(ctx, message.LoanID)
		if err != nil {
			logrus.Errorf("Failed to get loan by ID %d: %v", message.LoanID, err)
			resultCh <- AsyncResultData{} // Mengirimkan hasil kosong jika ada error
		} else {
			resultCh <- AsyncResultData{Loan: res} // Mengirimkan hasil sukses
		}
	}()

	go func() {
		defer wg.Done() // Menandakan bahwa goroutine ini selesai
		res, err := s.Repo.GetByLoanOrderNumber(ctx, message.LoanOrderNumber)
		if err != nil {
			logrus.Errorf("Failed to get loan funding by LoanOrderNumber %s: %v", message.LoanOrderNumber, err)
			resultCh <- AsyncResultData{} // Mengirimkan hasil kosong jika ada error
		} else {
			resultCh <- AsyncResultData{LoanFunding: res} // Mengirimkan hasil sukses
		}
	}()

	// Tunggu semua goroutine selesai
	wg.Wait()
	close(resultCh) // Menutup channel setelah semua goroutine selesai

	var loanFunding *repo.LoanFunding
	var loan *repo.Loan
	for res := range resultCh {
		if res.LoanFunding != nil {
			loanFunding = res.LoanFunding
		}
		if res.Loan != nil {
			loan = res.Loan
		}
	}

	// Mulai transaksi
	txnCtx := dbtxn.Begin(&ctx)

	defer func() {
		// Commit atau Rollback transaksi
		err := txnCtx.Commit()
		if err != nil {
			logrus.Errorf("Failed to commit transaction for LoanID %d: %v", loan.ID, err)
		}
	}()

	// Cek jika loanFunding atau loan tidak ditemukan
	if loanFunding == nil {
		logrus.Errorf("Loan funding process failed for LoanID %d: LoanFunding not found", message.LoanID)
		return fmt.Errorf("loan funding process failed for ID %d", message.LoanID)
	}
	if loan == nil {
		logrus.Errorf("Loan process failed for LoanID %d: Loan not found", message.LoanID)
		return fmt.Errorf("loan process failed for ID %d", message.LoanID)
	}

	isEligible := false

	defer func() {
		if isEligible {
			// Update loan invested amount
			loan.TotalInvestedAmount = loan.TotalInvestedAmount + loanFunding.InvestmentAmount
			loan.InvestorCount += int64(1)
			loan.UpdatedAt = time.Now()

			// check is principal amount is match with invested amount
			if loan.RequestAmount == loan.TotalInvestedAmount {
				loan.LoanStatus = enum.Invested
			}

			err := s.LoanRepo.Update(ctx, loan)
			if err != nil {
				txnCtx.AppendError(err)
				logrus.Errorf("Failed to update loan for LoanID %d: %v", loan.ID, err)
			}

			// Update loan funding
			loanFunding.Rate = loan.InvestmentPercentage
			loanFunding.Interest = utils.CalculateInterest(loanFunding.InvestmentAmount, loanFunding.Rate, loan.Tenures)
			loanFunding.ROI = loanFunding.InvestmentAmount + loanFunding.Interest
			loanFunding.Status = enum.LoanFundingInvested
			loanFunding.UpdatedAt = time.Now()

			err = s.Repo.Update(ctx, loanFunding)
			if err != nil {
				txnCtx.AppendError(err)
				logrus.Errorf("Failed to update loan funding for LoanOrderNumber %s: %v", loanFunding.LoanOrderNumber, err)
			}

			if loan.LoanStatus == enum.Invested {
				// init disburse
				disburseRequest := dto.LoanDisbursementRequestDTO{
					LoanID:         loan.ID,
					DisburseAmount: loan.TotalInvestedAmount,
				}
				err = s.DisburseSvc.Create(ctx, &disburseRequest)
				if err != nil {
					txnCtx.AppendError(err)
					logrus.Errorf("Failed while init disburse %s: %v", loanFunding.LoanOrderNumber, err)
				}
			}

			// todo trigger email
			email := SendEmailInput{
				To:      []string{loanFunding.LenderEmail},
				Subject: "aggrement",
				Body:    "hello world",
			}
			_ = s.MailSvc.SendEmail(ctx, email)

		} else {
			// update loan funding to failed
			loanFunding.Status = enum.LoanFundingFailed
			loanFunding.UpdatedAt = time.Now()
			err := s.Repo.Update(ctx, loanFunding)
			if err != nil {
				txnCtx.AppendError(err)
				logrus.Errorf("Failed to update loan funding for LoanOrderNumber %s: %v", loanFunding.LoanOrderNumber, err)
			}
		}
	}()

	// Cek status funding
	if loanFunding.Status != enum.LoanFundingPending {
		logrus.Infof("Loan funding for LoanID %d is not pending", message.LoanID)
		return nil
	}

	// Cek status loan
	if loan.LoanStatus != enum.Approved {
		logrus.Infof("Loan %d status is not approved, skipping funding process", loan.ID)
		return nil
	}

	// Cek funding deadline
	if loan.FundingDeadline.Before(time.Now()) {
		logrus.Warnf("Loan funding deadline for LoanID %d has passed", loan.ID)
		return errors.New("loan already expired")
	}

	// Cek apakah total pembayaran melebihi jumlah yang diminta
	if (loan.TotalInvestedAmount + loanFunding.InvestmentAmount) > loan.RequestAmount {
		logrus.Warnf("Total repayment amount exceeds requested loan amount for LoanID %d", loan.ID)
		return nil
	}

	// passed all validation
	isEligible = true

	logrus.Infof("Funding process completed successfully for LoanID %d", loan.ID)
	return nil
}

func (b *LoanFundingSvcImpl) validateRequestLoanFunding(loanFunding *dto.LoanFundingRequestDTO) string {
	// First validate using govalidator (it validates required and email fields)
	ok, err := govalidator.ValidateStruct(loanFunding)
	if !ok {
		return fmt.Sprintf("Validation failed: %s", err)
	}

	// Ensure that OrderNumber is not empty
	if loanFunding.OrderNumber == "" {
		return "OrderNumber must be provided"
	}

	// Ensure that LoanID is provided and greater than zero
	if loanFunding.LoanID <= 0 {
		return "LoanID must be provided and greater than zero"
	}

	// Ensure that LenderID is provided and greater than zero
	if loanFunding.LenderID <= 0 {
		return "LenderID must be provided and greater than zero"
	}

	// Ensure that InvestmentAmount is greater than zero
	if loanFunding.InvestmentAmount <= 0 {
		return "InvestmentAmount must be greater than zero"
	}

	// Validate LenderAgreementURL if it's provided (URL validation)
	if loanFunding.LenderAgreementURL != "" && !govalidator.IsURL(loanFunding.LenderAgreementURL) {
		return "LenderAgreementURL is invalid"
	}

	return ""
}

func (b *LoanFundingSvcImpl) publishFundingProcess(ctx context.Context, funding repo.LoanFunding) error {

	req := message2.FundingProcessMessage{
		LoanID:          funding.LoanID,
		LoanOrderNumber: funding.LoanOrderNumber,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		logrus.Error("Failed to marshal JSON:", err)
		return err
	}

	topic := string(consts.FundingProcessTopic)

	message := kafka.Message{
		Topic: topic,
		Value: jsonData,
	}

	err = b.KafkaWriter.WriteMessages(ctx, message)
	if err != nil {
		return err
	}

	return nil
}
