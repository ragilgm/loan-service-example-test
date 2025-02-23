package service

import (
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/test/loan-service/internal/dto"
	"github.com/test/loan-service/internal/dto/message"
	"github.com/test/loan-service/internal/enum"
	repo "github.com/test/loan-service/internal/repository"
	"github.com/test/loan-service/internal/utils"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
	"go.uber.org/dig"
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=$PROJ/internal/generated/mock/mock_$GOPACKAGE/$GOFILE

type (
	LoanRequest struct {
		Page   uint64
		Size   uint64
		Status *enum.LoanStatus
	}

	LoanSvc interface {
		Create(context.Context, *dto.LoanRequestDTO) (int64, error)
		ApprovalLoan(ctx context.Context, request message.UpdateLoanMessage) error
		DisburseLoan(ctx context.Context, request message.UpdateLoanMessage) error
		GetByID(ctx context.Context, loanID int64) (*dto.LoanResponseDTO, error)
		GetAllPage(ctx context.Context, request LoanRequest) ([]dto.LoanResponseDTO, int, error)
	}

	LoanSvcImpl struct {
		dig.In
		Repo            repo.LoanRepo
		LoanFundingRepo repo.LoanFundingRepo
		LoanDetailSvc   LoanDetailSvc
		LoanApprovalSvc LoanApprovalSvc
	}
)

func NewLoanSvc(impl LoanSvcImpl) LoanSvc {
	return &impl
}

func (b *LoanSvcImpl) Create(ctx context.Context, loanRequest *dto.LoanRequestDTO) (int64, error) {
	// Validate request
	if errMsg := b.validateRequestLoan(loanRequest); errMsg != "" {
		log.WithFields(log.Fields{
			"borrowerID":    loanRequest.BorrowerID,
			"requestAmount": loanRequest.RequestAmount,
		}).Errorf("Validation failed: %s", errMsg)
		return -1, errors.New(errMsg)
	}

	var loan repo.Loan
	err := mapstructure.Decode(loanRequest, &loan)
	if err != nil {
		log.WithError(err).Error("Failed to decode loan request")
		return -1, errors.New("99999")
	}

	// Generate loan code
	loanCode := utils.GenerateAlphanumericCode(10)
	loan.LoanCode = loanCode

	loan.InvestmentPercentage = loan.Rate

	// Set initial loan status
	loan.LoanStatus = enum.Proposed

	// Start transactional
	txnCtx := dbtxn.Begin(&ctx)
	defer func() {
		if err := txnCtx.Commit(); err != nil {
			log.WithError(err).Error("Transaction commit failed")
			txnCtx.AppendError(err)
		}
	}()

	// Create loan in repository
	id, err := b.Repo.Create(ctx, &loan)
	if err != nil {
		log.WithFields(log.Fields{
			"loanCode": loanCode,
		}).WithError(err).Error("Failed to create loan in repo")
		txnCtx.AppendError(err)
		return -1, errors.New("99999")
	}

	// Create loan detail
	_, err = b.createLoanDetail(ctx, loanRequest, id)
	if err != nil {
		log.WithFields(log.Fields{
			"loanID": id,
		}).WithError(err).Error("Failed to create loan detail")
		txnCtx.AppendError(err)
		return -1, errors.New("99999")
	}

	// Create initial approval
	_, err = b.createInitialApproval(ctx, id)
	if err != nil {
		log.WithFields(log.Fields{
			"loanID": id,
		}).WithError(err).Error("Failed to create initial approval")
		txnCtx.AppendError(err)
		return -1, errors.New("99999")
	}

	log.WithFields(log.Fields{
		"loanID":   id,
		"loanCode": loanCode,
	}).Info("Loan created successfully")
	return id, nil
}

func (b *LoanSvcImpl) ApprovalLoan(ctx context.Context, request message.UpdateLoanMessage) error {
	// Retrieve loan based on ID
	loan, err := b.Repo.GetByID(ctx, request.LoanID)
	if err != nil {
		log.WithFields(log.Fields{
			"loanID": request.LoanID,
		}).WithError(err).Error("Failed to retrieve loan from repo")
		return errors.New("99999")
	}

	if loan == nil {
		log.WithFields(log.Fields{
			"loanID": request.LoanID,
		}).Warn("Loan not found")
		return errors.New("99999")
	}

	// Validate status transition
	isValid := b.isValidStatusTransition(loan.LoanStatus, request.LoanStatus)
	if !isValid {
		log.WithFields(log.Fields{
			"currentStatus": loan.LoanStatus,
			"newStatus":     request.LoanStatus,
		}).Error("Invalid status transition")
		return errors.New("99999")
	}

	// Update loan status
	loan.LoanStatus = request.LoanStatus
	if loan.LoanStatus == enum.Approved {
		fundingDeadline := time.Now().Add(7 * 24 * time.Hour) // 7 days from now
		loan.FundingDeadline = &fundingDeadline
	}

	loan.UpdatedAt = time.Now()

	// Start transaction to update loan
	txnCtx := dbtxn.Begin(&ctx)
	defer func() {
		if err = txnCtx.Commit(); err != nil {
			log.WithError(err).Error("Transaction commit failed")
			txnCtx.AppendError(err)
		}
	}()

	err = b.Repo.Update(ctx, loan)
	if err != nil {
		log.WithFields(log.Fields{
			"loanID": request.LoanID,
		}).WithError(err).Error("Failed to update loan")
		txnCtx.AppendError(err)
		return errors.New("99999")
	}

	log.WithFields(log.Fields{
		"loanID":    request.LoanID,
		"newStatus": loan.LoanStatus,
	}).Info("Loan updated successfully")
	return nil
}

func (b *LoanSvcImpl) DisburseLoan(ctx context.Context, request message.UpdateLoanMessage) error {
	// Retrieve loan based on ID
	loan, err := b.Repo.GetByID(ctx, request.LoanID)
	if err != nil {
		log.WithFields(log.Fields{
			"loanID": request.LoanID,
		}).WithError(err).Error("Failed to retrieve loan from repo")
		return errors.New("99999")
	}

	if loan == nil {
		log.WithFields(log.Fields{
			"loanID": request.LoanID,
		}).Warn("Loan not found")
		return errors.New("99999")
	}

	if request.LoanStatus != enum.Disbursed {
		log.Errorf("Loan status %s is not disbursed", request.LoanStatus)
		return nil
	}

	// Validate status transition
	isValid := b.isValidStatusTransition(loan.LoanStatus, request.LoanStatus)
	if !isValid {
		log.WithFields(log.Fields{
			"currentStatus": loan.LoanStatus,
			"newStatus":     request.LoanStatus,
		}).Error("Invalid status transition")
		return errors.New("99999")
	}

	// change status
	loan.LoanStatus = request.LoanStatus
	// calculate interest for borrower
	loan.TotalInterest = utils.CalculateInterest(loan.RequestAmount, loan.Rate, loan.Tenures)
	// calculate total repayment amount ( total amount which nedd borrower pay )
	loan.TotalRepaymentAmount = loan.RequestAmount + loan.TotalInterest
	loan.UpdatedAt = time.Now()

	// Start transaction to update loan
	txnCtx := dbtxn.Begin(&ctx)
	defer func() {
		if err = txnCtx.Commit(); err != nil {
			log.WithError(err).Error("Transaction commit failed")
			txnCtx.AppendError(err)
		}
	}()

	err = b.Repo.Update(ctx, loan)
	if err != nil {
		log.WithFields(log.Fields{
			"loanID": request.LoanID,
		}).WithError(err).Error("Failed to update loan")
		txnCtx.AppendError(err)
		return errors.New("99999")
	}

	// update loan funding
	var loanFunding []repo.LoanFunding
	loanFunding, err = b.LoanFundingRepo.GetByLoanID(ctx, loan.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"loanID": request.LoanID,
		}).WithError(err).Error("Failed to scan loan funding")
		txnCtx.AppendError(err)
		return errors.New("99999")
	}

	for _, funding := range loanFunding {
		// change funding status from invested to on going , becuase the loan already disbursed
		funding.Status = enum.LoanFundingOngoing
		funding.UpdatedAt = time.Now()
		err = b.LoanFundingRepo.Update(ctx, &funding)
		if err != nil {
			log.WithFields(log.Fields{
				"loanID": request.LoanID,
			}).WithError(err).Error("Failed to update loan funding")
			txnCtx.AppendError(err)
			return errors.New("99999")
		}
	}

	// TODO : generate repayment schedule borrower

	log.WithFields(log.Fields{
		"loanID":    request.LoanID,
		"newStatus": loan.LoanStatus,
	}).Info("Loan updated successfully")
	return nil
}

func (b *LoanSvcImpl) GetByID(ctx context.Context, loanID int64) (*dto.LoanResponseDTO, error) {
	// Log request to get loan by ID
	log.WithFields(log.Fields{
		"loanID": loanID,
	}).Info("Fetching loan by ID")

	// Retrieve the loan using the repository
	loan, err := b.Repo.GetByID(ctx, loanID)
	if err != nil {
		log.WithFields(log.Fields{
			"loanID": loanID,
		}).WithError(err).Error("Failed to retrieve loan from repository")
		return nil, errors.New("99999")
	}

	// If loan is not found, return an error
	if loan == nil {
		log.WithFields(log.Fields{
			"loanID": loanID,
		}).Warn("Loan not found")
		return nil, errors.New("10001")
	}

	// Log successful retrieval of loan
	log.WithFields(log.Fields{
		"loanID":   loanID,
		"loanCode": loan.LoanCode,
	}).Info("Loan retrieved successfully")

	// Map the loan to the response DTO (data transfer object)
	var loanResponse dto.LoanResponseDTO
	err = mapstructure.Decode(loan, &loanResponse)
	if err != nil {
		log.WithFields(log.Fields{
			"loanID": loanID,
		}).WithError(err).Error("Failed to map loan to DTO")
		return nil, errors.New("99999")
	}

	// Set additional fields in loan response DTO
	loanResponse.FundingDeadline = loan.FundingDeadline
	loanResponse.CreatedAt = loan.CreatedAt
	loanResponse.UpdatedAt = loan.UpdatedAt
	loanResponse.DeletedAt = loan.DeletedAt

	// Log successful mapping to DTO
	log.WithFields(log.Fields{
		"loanID":   loanID,
		"loanCode": loan.LoanCode,
	}).Info("Loan response DTO created successfully")

	// Return the loan response DTO
	return &loanResponse, nil
}

func (b *LoanSvcImpl) createLoanDetail(ctx context.Context, loan *dto.LoanRequestDTO, loanID int64) (int64, error) {
	log.WithFields(log.Fields{
		"loanID": loanID,
	}).Info("Creating loan detail")

	// Create loan detail
	detailRequest := LoanDetailRequest{
		BorrowerID: loan.BorrowerID,
		LoanID:     loanID,
		Detail:     &loan.Detail,
	}

	id, err := b.LoanDetailSvc.Create(ctx, &detailRequest)
	if err != nil {
		log.WithFields(log.Fields{
			"loanID": loanID,
		}).WithError(err).Error("Failed to create loan detail")
		return -1, errors.New("99999")
	}

	log.WithFields(log.Fields{
		"loanID":   loanID,
		"detailID": id,
	}).Info("Loan detail created successfully")
	return id, nil
}

func (b *LoanSvcImpl) createInitialApproval(ctx context.Context, loanID int64) (int64, error) {
	log.WithFields(log.Fields{
		"loanID": loanID,
	}).Info("Creating initial approval")

	// Create initial approval
	approvalRequest := dto.LoanApprovalRequestDTO{
		LoanID: loanID,
	}

	id, err := b.LoanApprovalSvc.Create(ctx, &approvalRequest)
	if err != nil {
		log.WithFields(log.Fields{
			"loanID": loanID,
		}).WithError(err).Error("Failed to create initial approval")
		return -1, errors.New("99999")
	}

	log.WithFields(log.Fields{
		"loanID":     loanID,
		"approvalID": id,
	}).Info("Initial approval created successfully")
	return id, nil
}

func (b *LoanSvcImpl) validateRequestLoan(loan *dto.LoanRequestDTO) string {
	// Validate the loan request
	ok, err := govalidator.ValidateStruct(loan)
	if !ok {
		log.WithError(err).Error("Loan validation failed")
		return "10003"
	}

	// Additional loan-specific validations
	if loan.BorrowerID <= 0 {
		log.Error("BorrowerID must be provided and greater than zero")
		return "10003"
	}

	if loan.RequestAmount <= 0 {
		log.Error("RequestAmount must be greater than zero")
		return "10003"
	}

	if loan.LoanGrade == "" {
		log.Error("LoanGrade must be provided")
		return "10003"
	}

	if !loan.LoanType.IsValid() {
		log.Error("Invalid LoanType")
		return "10003"

	}

	if loan.Rate < 0 {
		log.Error("Rate must be non-negative")
		return "10003"
	}

	if loan.Tenures <= 0 {
		log.Error("Tenures must be greater than zero")
		return "10003"
	}

	return ""
}

func (b *LoanSvcImpl) isValidStatusTransition(currentStatus, newStatus enum.LoanStatus) bool {
	// Valid status transitions for loans
	validTransitions := map[enum.LoanStatus][]enum.LoanStatus{
		enum.Proposed:  {enum.Approved, enum.Rejected},
		enum.Approved:  {enum.Invested},
		enum.Invested:  {enum.Disbursed},
		enum.Disbursed: {enum.Completed},
	}

	allowedStatuses, ok := validTransitions[currentStatus]
	if !ok {
		return false
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			return true
		}
	}

	return false
}

func (b *LoanSvcImpl) GetAllPage(ctx context.Context, request LoanRequest) ([]dto.LoanResponseDTO, int, error) {
	// Log the request details, including pagination parameters
	log.WithFields(log.Fields{
		"page": request.Page,
		"size": request.Size,
	}).Info("Fetching paginated loan data")

	// Calculate offset based on page and size
	offset := (request.Page - 1) * request.Size
	var repoReq repo.LoanRequest
	err := mapstructure.Decode(request, &repoReq)
	if err != nil {
		log.WithError(err).Error("Failed to decode request to repository format")
		return nil, 0, errors.New("99999")
	}
	repoReq.Offset = offset

	// Fetch loans from the repository
	loans, totalRecords, err := b.Repo.GetAllPage(ctx, repoReq)
	if err != nil {
		log.WithError(err).Error("Failed to fetch loans from repository")
		return nil, 0, errors.New("99999")
	}

	// Map loans from repository model to response DTO
	var loanDTOs []dto.LoanResponseDTO
	for _, loan := range loans {
		var loanDTO dto.LoanResponseDTO
		err = mapstructure.Decode(loan, &loanDTO)
		if err != nil {
			log.WithError(err).WithField("loanID", loan.ID).Error("Failed to map loan to response DTO")
			return nil, 0, errors.New("99999")
		}
		loanDTO.CreatedAt = loan.CreatedAt
		loanDTO.UpdatedAt = loan.UpdatedAt
		loanDTO.DeletedAt = loan.DeletedAt

		var loanDetail *dto.LoanDetailResponseDTO
		loanDetail, err = b.LoanDetailSvc.GetByLoanID(ctx, loan.ID)
		if err != nil {
			log.WithError(err).WithField("loanID", loan.ID).Error("Failed to get loan detail")
			return nil, 0, errors.New("99999")
		}
		if loanDetail != nil {
			loanDTO.LoanDetail = loanDetail
		}

		loanDTOs = append(loanDTOs, loanDTO)
	}

	// Log success with the total records found
	log.WithFields(log.Fields{
		"totalRecords": totalRecords,
	}).Info("Successfully fetched paginated loans")

	return loanDTOs, int(totalRecords), nil
}
