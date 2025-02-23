package service

import (
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/test/loan-service/internal/dto"
	repo "github.com/test/loan-service/internal/repository"
	"go.uber.org/dig"
)

type (
	LoanDetailRequest struct {
		Detail     *dto.LoanDetailRequestDTO
		BorrowerID int64
		LoanID     int64
	}

	LoanDetailSvc interface {
		Create(context.Context, *LoanDetailRequest) (int64, error)
		GetByLoanID(ctx context.Context, loanID int64) (*dto.LoanDetailResponseDTO, error)
	}

	LoanDetailSvcImpl struct {
		dig.In
		Repo repo.LoanDetailRepo
	}
)

func NewLoanDetailSvc(impl LoanDetailSvcImpl) LoanDetailSvc {
	return &impl
}

func (b *LoanDetailSvcImpl) Create(ctx context.Context, detailRequest *LoanDetailRequest) (int64, error) {
	logrus.WithFields(logrus.Fields{
		"loanID":     detailRequest.LoanID,
		"borrowerID": detailRequest.BorrowerID,
	}).Info("Creating loan detail")

	if errMsg := b.validateLoanDetail(detailRequest); errMsg != "" {
		logrus.WithField("loanID", detailRequest.LoanID).Error("Loan detail validation failed: " + errMsg)
		return -1, errors.New(errMsg)
	}

	var loanDetail repo.LoanDetail
	err := mapstructure.Decode(detailRequest.Detail, &loanDetail)
	if err != nil {
		logrus.WithField("loanID", detailRequest.LoanID).WithError(err).Error("Failed to map request to loan detail")
		return -1, errors.New("99999")
	}

	loanDetail.LoanID = detailRequest.LoanID
	loanDetail.BorrowerID = detailRequest.BorrowerID

	id, err := b.Repo.Create(ctx, &loanDetail)
	if err != nil {
		logrus.WithField("loanID", detailRequest.LoanID).WithError(err).Error("Failed to create loan detail")
		return -1, errors.New("99999")
	}

	logrus.WithField("loanDetailID", id).Info("Loan detail created successfully")
	return id, nil
}

func (b *LoanDetailSvcImpl) GetByLoanID(ctx context.Context, loanID int64) (*dto.LoanDetailResponseDTO, error) {
	logrus.WithField("loanID", loanID).Info("Fetching loan details")

	loanDetail, err := b.Repo.GetByLoanID(ctx, loanID)
	if err != nil {
		logrus.WithField("loanID", loanID).WithError(err).Error("Failed to get loan details")
		return nil, errors.New("99999")
	}

	if loanDetail == nil {
		logrus.WithField("loanID", loanID).Warn("Loan details not found")
		return nil, errors.New("99999")
	}

	var detailRes dto.LoanDetailResponseDTO
	err = mapstructure.Decode(loanDetail, &detailRes)
	if err != nil {
		logrus.WithField("loanID", loanID).WithError(err).Error("Failed to decode loan details")
		return nil, errors.New("99999")
	}

	detailRes.CreatedAt = loanDetail.CreatedAt
	detailRes.UpdatedAt = loanDetail.UpdatedAt
	detailRes.DeletedAt = loanDetail.DeletedAt

	logrus.WithField("loanID", loanID).Info("Loan details fetched successfully")
	return &detailRes, nil
}

func (b *LoanDetailSvcImpl) validateLoanDetail(loanDetail *LoanDetailRequest) string {
	ok, err := govalidator.ValidateStruct(loanDetail)
	if !ok || err != nil {
		logrus.WithFields(logrus.Fields{
			"loanID":     loanDetail.LoanID,
			"borrowerID": loanDetail.BorrowerID,
		}).Error("Loan detail validation failed")
		return "10003"
	}

	logrus.WithFields(logrus.Fields{
		"loanID":     loanDetail.LoanID,
		"borrowerID": loanDetail.BorrowerID,
	}).Info("Loan detail validated successfully")
	return ""
}
