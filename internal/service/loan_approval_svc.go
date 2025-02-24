package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
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

//go:generate mockgen -source=$GOFILE -destination=$PROJ/internal/generated/mock/mock_$GOPACKAGE/$GOFILE

type (
	LoanApprovalSvc interface {
		Create(context.Context, *dto.LoanApprovalRequestDTO) (int64, error)
		Update(context.Context, int64, *dto.UpdateLoanApprovalRequestDTO) error
		GetAllPage(ctx context.Context, request models.LoanApprovalRequest) ([]dto.LoanApprovalResponseDTO, int, error)
	}

	LoanApprovalSvcImpl struct {
		dig.In
		Repo                 repo.LoanApprovalRepo
		ApprovalDocumentRepo repo.ApprovalDocumentRepo
		KafkaWriter          *kafka.Writer
		Validator            validator.LoanApprovalValidatorImpl
	}
)

func NewLoanApprovalSvc(impl LoanApprovalSvcImpl) LoanApprovalSvc {
	return &impl
}

func (b *LoanApprovalSvcImpl) Create(ctx context.Context, loanRequest *dto.LoanApprovalRequestDTO) (int64, error) {
	logrus.Infof("Creating loan approval for request: %+v", loanRequest)

	// validate request
	err := b.Validator.ValidateCreate(loanRequest)

	if err != nil {
		logrus.Errorf("Validation failed for loan approval: %s", err)
		return -1, err
	}

	var approval repo.LoanApproval
	err = mapstructure.Decode(loanRequest, &approval)
	if err != nil {
		logrus.Errorf("Error decoding loan request to loan approval model: %v", err)
		return -1, errors.New("99999")
	}

	// generate approval number
	approval.ApprovalNumber = utils.GenerateAlphanumericCode(10)

	// initial status
	approval.ApprovalStatus = enum.ApprovalPending
	approval.CreatedAt = time.Now()
	approval.UpdatedAt = time.Now()

	id, err := b.Repo.Create(ctx, &approval)
	if err != nil {
		logrus.Errorf("Error creating loan approval in repository: %v", err)
		return -1, errors.New("99999")
	}

	return id, nil
}

func (b *LoanApprovalSvcImpl) Update(ctx context.Context, approvalId int64, requestDTO *dto.UpdateLoanApprovalRequestDTO) error {
	logrus.Infof("Updating loan approval with ID: %d", approvalId)

	// validate request
	err := b.Validator.ValidateUpdate(requestDTO)
	if err != nil {
		logrus.Errorf("Validation failed for loan approval update: %s", err)
		return err
	}

	approval, err := b.Repo.GetByID(ctx, approvalId)
	if err != nil {
		logrus.Errorf("Error fetching loan approval with ID: %d", approvalId)
		return errors.New("99999")
	}
	if approval == nil {
		logrus.Warnf("Loan approval not found with ID: %d", approvalId)
		return errors.New("10001")
	}

	// validate transition status
	if !b.Validator.ValidateTransitionStatus(approval.ApprovalStatus, requestDTO.ApprovalStatus) {
		logrus.Warnf("Invalid status transition from %s to %s", approval.ApprovalStatus, requestDTO.ApprovalStatus)
		return errors.New("10003")
	}

	// start transactional
	txnCtx := dbtxn.Begin(&ctx)
	defer func() {
		// Commit or Rollback transactional
		if err := txnCtx.Commit(); err != nil {
			logrus.Errorf("Error committing transaction: %v", err)
		}
	}()

	approvalDate := time.Now()
	approval.StaffID = &requestDTO.StaffID
	approval.ApprovalStatus = requestDTO.ApprovalStatus
	approval.ApprovalDate = &approvalDate
	approval.UpdatedAt = approvalDate

	// update loan approval
	logrus.Infof("Updating loan approval with ID: %d in repository", approvalId)
	err = b.Repo.Update(ctx, approval)
	if err != nil {
		logrus.Errorf("Error updating loan approval with ID: %d: %v", approvalId, err)
		txnCtx.AppendError(err)
		return errors.New("99999")
	}

	for _, document := range requestDTO.ApprovalDocuments {
		logrus.Infof("Processing document for loan approval ID: %d", approvalId)
		var approvalDocument repo.ApprovalDocument
		err = mapstructure.Decode(document, &approvalDocument)
		if err != nil {
			logrus.Errorf("Error decoding document: %v", err)
			txnCtx.AppendError(err)
			return errors.New("99999")
		}

		now := time.Now()
		approvalDocument.LoanApprovalID = approval.ID
		approvalDocument.CreatedAt = now
		approvalDocument.UpdatedAt = now

		_, err = b.ApprovalDocumentRepo.Create(ctx, &approvalDocument)
		if err != nil {
			logrus.Errorf("Error saving approval document for loan approval ID: %d: %v", approvalId, err)
			txnCtx.AppendError(err)
			return errors.New("99999")
		}
	}

	// publish kafka loan update
	logrus.Infof("Publishing loan update to Kafka for approval ID: %d", approvalId)
	err = b.publishLoanApproval(ctx, approvalId, approval.ApprovalStatus, err)
	if err != nil {
		txnCtx.AppendError(err)
		return errors.New("99999")
	}

	logrus.Infof("Successfully updated loan approval with ID: %d", approvalId)
	return nil
}

func (b *LoanApprovalSvcImpl) publishLoanApproval(ctx context.Context, approvalId int64, approvalStatus enum.ApprovalStatus, err error) error {
	// publish kafka for update loan
	req := message2.UpdateLoanMessage{
		LoanID:     approvalId,
		LoanStatus: enum.LoanStatus(approvalStatus),
	}

	logrus.Infof("Marshalling loan update message for approval ID: %d", approvalId)
	jsonData, err := json.Marshal(req)
	if err != nil {
		logrus.Errorf("Failed to marshal loan update message for approval ID: %d: %v", approvalId, err)
		return errors.New("99999")
	}

	topic := string(consts.ApprovalLoanTopic)
	message := kafka.Message{
		Topic: topic,
		Value: jsonData,
	}

	logrus.Infof("Sending loan update message to Kafka topic: %s", topic)
	err = b.KafkaWriter.WriteMessages(ctx, message)
	if err != nil {
		logrus.Errorf("Failed to send loan update message to Kafka topic: %s: %v", topic, err)
		return errors.New("99999")
	}

	logrus.Infof("Successfully published loan update for approval ID: %d", approvalId)
	return nil
}

func (b *LoanApprovalSvcImpl) GetAllPage(ctx context.Context, request models.LoanApprovalRequest) ([]dto.LoanApprovalResponseDTO, int, error) {
	logrus.Infof("Fetching all loan approvals with pagination: Page: %d, Size: %d", request.Page, request.Size)

	// Calculate offset
	offset := (request.Page - 1) * request.Size

	var repoReq repo.LoanApprovalRequest
	err := mapstructure.Decode(request, &repoReq)
	if err != nil {
		logrus.Errorf("Error decoding loan approval request: %v", err)
		return nil, 0, errors.New("99999")
	}
	repoReq.Offset = offset

	// Get loans from repo
	approvals, totalRecords, err := b.Repo.GetAllPage(ctx, repoReq)
	if err != nil {
		logrus.Errorf("Error fetching paginated loan approvals: %v", err)
		return nil, 0, errors.New("99999")
	}

	// Convert to DTO
	approvalDTOs := []dto.LoanApprovalResponseDTO{}
	for _, approval := range approvals {
		var approvalRes dto.LoanApprovalResponseDTO
		err = mapstructure.Decode(approval, &approvalRes)
		if err != nil {
			logrus.Errorf("Error decoding approval to response DTO: %v", err)
			return nil, 0, errors.New("99999")
		}

		approvalRes.CreatedAt = approval.CreatedAt
		approvalRes.UpdatedAt = approval.UpdatedAt
		approvalRes.DeletedAt = approval.DeletedAt

		// Fetch documents if approval is approved
		if approvalRes.ApprovalStatus == enum.ApprovalApproved {
			logrus.Infof("Fetching documents for approved loan ID: %d", approval.ID)
			var approvalDocs []repo.ApprovalDocument
			approvalDocs, err = b.ApprovalDocumentRepo.GetByApprovalID(ctx, approval.ID)
			if err != nil {
				logrus.Errorf("Error fetching documents for loan ID: %d: %v", approval.ID, err)
				return nil, 0, errors.New("99999")
			}

			var approvalDocDTOs []dto.ApprovalDocumentResponseDTO
			for _, approvalDoc := range approvalDocs {
				var apprDocDto dto.ApprovalDocumentResponseDTO
				err = mapstructure.Decode(approvalDoc, &apprDocDto)
				if err != nil {
					logrus.Errorf("Error decoding document to DTO: %v", err)
					return nil, 0, errors.New("99999")
				}
				apprDocDto.CreatedAt = approvalDoc.CreatedAt
				apprDocDto.UpdatedAt = approvalDoc.UpdatedAt
				apprDocDto.DeletedAt = approvalDoc.DeletedAt

				approvalDocDTOs = append(approvalDocDTOs, apprDocDto)
			}
			approvalRes.ApprovalDocument = &approvalDocDTOs
		}

		approvalDTOs = append(approvalDTOs, approvalRes)
	}

	logrus.Infof("Successfully fetched %d loan approvals", len(approvalDTOs))
	return approvalDTOs, int(totalRecords), nil
}
