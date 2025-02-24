package validator

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/test/loan-service/internal/dto"
	"go.uber.org/dig"
)

type LoanFundingValidatorImpl struct {
	dig.In
}

func NewLoanFundingValidator(impl LoanFundingValidatorImpl) CustomValidator {
	return &impl
}

func (l LoanFundingValidatorImpl) ValidateCreate(data interface{}) error {

	var loanFunding dto.LoanFundingRequestDTO
	err := mapstructure.Decode(data, &loanFunding)
	if err != nil {
		return errors.New("99999")
	}

	ok, err := govalidator.ValidateStruct(loanFunding)
	if !ok {
		log.Errorf("Validation failed: %s", err)
		return errors.New("10003")
	}

	// Ensure that OrderNumber is not empty
	if loanFunding.OrderNumber == "" {
		log.Errorf("OrderNumber must be provided")
		return errors.New("10003")
	}

	// Ensure that LoanID is provided and greater than zero
	if loanFunding.LoanID <= 0 {
		log.Errorf("LoanID must be provided and greater than zero")
		return errors.New("10003")
	}

	// Ensure that LenderID is provided and greater than zero
	if loanFunding.LenderID <= 0 {
		log.Errorf("LenderID must be provided and greater than zero")
		return errors.New("10003")
	}

	// Ensure that InvestmentAmount is greater than zero
	if loanFunding.InvestmentAmount <= 0 {
		log.Errorf("InvestmentAmount must be greater than zero")
		return errors.New("10003")
	}

	// Validate LenderAgreementURL if it's provided (URL validation)
	if loanFunding.LenderAgreementURL != "" && !govalidator.IsURL(loanFunding.LenderAgreementURL) {
		log.Errorf("LenderAgreementURL is invalid")
		return errors.New("10003")
	}

	return nil

}

func (l LoanFundingValidatorImpl) ValidateUpdate(data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (l LoanFundingValidatorImpl) ValidateTransitionStatus(from interface{}, to interface{}) bool {
	//TODO implement me
	panic("implement me")
}
