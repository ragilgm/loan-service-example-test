package validator

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/test/loan-service/internal/service/models"
	"go.uber.org/dig"
)

type LoanDetailValidatorImpl struct {
	dig.In
}

func NewLoanDetailValidator(impl LoanDetailValidatorImpl) CustomValidator {
	return &impl
}

func (l LoanDetailValidatorImpl) ValidateCreate(data interface{}) error {

	var loanDetail models.LoanDetailRequest
	err := mapstructure.Decode(data, &loanDetail)
	if err != nil {
		return errors.New("99999")
	}

	ok, err := govalidator.ValidateStruct(loanDetail)
	if !ok || err != nil {
		logrus.WithFields(logrus.Fields{
			"loanID":     loanDetail.LoanID,
			"borrowerID": loanDetail.BorrowerID,
		}).Error("Loan detail validation failed")
		return errors.New("10003")
	}

	logrus.WithFields(logrus.Fields{
		"loanID":     loanDetail.LoanID,
		"borrowerID": loanDetail.BorrowerID,
	}).Info("Loan detail validated successfully")
	return nil

}

func (l LoanDetailValidatorImpl) ValidateUpdate(data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (l LoanDetailValidatorImpl) ValidateTransitionStatus(from interface{}, to interface{}) bool {
	//TODO implement me
	panic("implement me")
}
