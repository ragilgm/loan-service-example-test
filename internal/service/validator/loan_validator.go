package validator

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/test/loan-service/internal/dto"
	"github.com/test/loan-service/internal/enum"
	"go.uber.org/dig"
)

type LoanValidatorImpl struct {
	dig.In
}

func NewLoanValidator(impl LoanValidatorImpl) CustomValidator {
	return &impl
}

func (lv *LoanValidatorImpl) ValidateCreate(data interface{}) error {
	// Validate the loan request
	var loan dto.LoanRequestDTO
	err := mapstructure.Decode(data, &loan)
	if err != nil {
		return errors.New("99999")
	}
	ok, err := govalidator.ValidateStruct(loan)
	if !ok {
		log.WithError(err).Error("Loan validation failed")
		return errors.New("10003")
	}

	// Additional loan-specific validations
	if loan.BorrowerID <= 0 {
		log.Error("BorrowerID must be provided and greater than zero")
		return errors.New("10003")
	}

	if loan.RequestAmount <= 0 {
		log.Error("RequestAmount must be greater than zero")
		return errors.New("10003")
	}

	if loan.LoanGrade == "" {
		log.Error("LoanGrade must be provided")
		return errors.New("10003")
	}

	if !loan.LoanType.IsValid() {
		log.Error("Invalid LoanType")
		return errors.New("10003")

	}

	if loan.Rate < 0 {
		log.Error("Rate must be non-negative")
		return errors.New("10003")
	}

	if loan.Tenures <= 0 {
		log.Error("Tenures must be greater than zero")
		return errors.New("10003")
	}

	return nil
}

func (lv *LoanValidatorImpl) ValidateUpdate(data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (lv *LoanValidatorImpl) ValidateTransitionStatus(from interface{}, to interface{}) bool {
	var currentStatus enum.LoanStatus
	currentStatus = from.(enum.LoanStatus)
	var changeStatus enum.LoanStatus
	changeStatus = to.(enum.LoanStatus)

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
		if status == changeStatus {
			return true
		}
	}

	return false
}
