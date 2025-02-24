package service

import (
	"context"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/test/loan-service/internal/dto"
	repo "github.com/test/loan-service/internal/repository"
	"github.com/test/loan-service/internal/service/models"
	"github.com/test/loan-service/internal/service/validator"
	"go.uber.org/dig"
)

type (
	LoanDetailSvc interface {
		Create(context.Context, *models.LoanDetailRequest) (int64, error)
		GetByLoanID(ctx context.Context, loanID int64) (*dto.LoanDetailResponseDTO, error)
	}

	LoanDetailSvcImpl struct {
		dig.In
		Repo      repo.LoanDetailRepo
		Validator validator.LoanDetailValidatorImpl
	}
)

func NewLoanDetailSvc(impl LoanDetailSvcImpl) LoanDetailSvc {
	return &impl
}

func (b *LoanDetailSvcImpl) Create(ctx context.Context, detailRequest *models.LoanDetailRequest) (int64, error) {
	logrus.WithFields(logrus.Fields{
		"loanID":     detailRequest.LoanID,
		"borrowerID": detailRequest.BorrowerID,
	}).Info("Creating loan detail")

	err := b.Validator.ValidateCreate(detailRequest)

	if err != nil {
		logrus.WithField("loanID", detailRequest.LoanID).Error("Loan detail validation failed: " + err.Error())
		return -1, err
	}

	var loanDetail repo.LoanDetail
	err = mapstructure.Decode(detailRequest.Detail, &loanDetail)
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
