package validator

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/test/loan-service/internal/dto"
	"github.com/test/loan-service/internal/enum"
	"go.uber.org/dig"
)

type LoanApprovalValidatorImpl struct {
	dig.In
}

func NewLoanApprovalValidator(impl LoanApprovalValidatorImpl) CustomValidator {
	return &impl
}

func (lv *LoanApprovalValidatorImpl) ValidateCreate(data interface{}) error {
	var loan dto.LoanApprovalRequestDTO
	err := mapstructure.Decode(data, &loan)
	if err != nil {
		return errors.New("99999")
	}

	ok, err := govalidator.ValidateStruct(loan)
	if !ok {
		logrus.Errorf("Validation failed for loan approval: %s", err)
		return errors.New("10003")
	}

	return nil

}

func (lv *LoanApprovalValidatorImpl) ValidateUpdate(data interface{}) error {

	var loan dto.UpdateLoanApprovalRequestDTO
	err := mapstructure.Decode(data, &loan)
	if err != nil {
		return errors.New("99999")
	}

	ok, err := govalidator.ValidateStruct(loan)
	if !ok {
		logrus.Errorf("Validation failed for loan approval update: %s", err)
		return errors.New("10003")
	}

	if !loan.ApprovalStatus.IsValid() {
		logrus.Errorf("Invalid approval status: %s", loan.ApprovalStatus)
		return errors.New("10003")
	}

	return nil
}

func (lv *LoanApprovalValidatorImpl) ValidateTransitionStatus(from interface{}, to interface{}) bool {

	validTransitions := map[enum.ApprovalStatus][]enum.ApprovalStatus{
		enum.ApprovalPending: {enum.ApprovalApproved, enum.ApprovalRejected},
	}

	var currentStatus enum.ApprovalStatus
	currentStatus = from.(enum.ApprovalStatus)
	var changeStatus enum.ApprovalStatus
	changeStatus = to.(enum.ApprovalStatus)

	allowedStatuses, ok := validTransitions[currentStatus]
	if !ok {
		logrus.Warnf("No valid transitions for current status: %s", currentStatus)
		return false
	}

	for _, status := range allowedStatuses {
		if status == changeStatus {
			return true
		}
	}

	logrus.Warnf("Invalid status transition from %s to %s", currentStatus, changeStatus)
	return false
}
