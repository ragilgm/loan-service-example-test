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

type LoanDisbursementValidatorImpl struct {
	dig.In
}

func NewLoanDisbursementValidator(impl LoanDisbursementValidatorImpl) CustomValidator {
	return &impl
}

func (l LoanDisbursementValidatorImpl) ValidateCreate(data interface{}) error {

	var disbursement dto.LoanDisbursementRequestDTO
	err := mapstructure.Decode(data, &disbursement)
	if err != nil {
		return errors.New("99999")
	}

	ok, err := govalidator.ValidateStruct(disbursement)
	if !ok {
		log.Printf("Validation failed: %s", err)
		return errors.New("10003")
	}

	if disbursement.DisburseAmount <= 0 {
		log.Printf("DisburseAmount must be greater than zero")
		return errors.New("10003")
	}

	if disbursement.LoanID <= 0 {
		log.Printf("LoanID must be provided and greater than zero")
		return errors.New("10003")
	}

	log.Printf("Validation passed for loan disbursement request: LoanID=%d", disbursement.LoanID)
	return nil
}

func (l LoanDisbursementValidatorImpl) ValidateUpdate(data interface{}) error {

	var disbursement dto.UpdateLoanDisbursementRequestDTO
	err := mapstructure.Decode(data, &disbursement)
	if err != nil {
		return errors.New("99999")
	}

	ok, err := govalidator.ValidateStruct(disbursement)
	if !ok {
		log.Printf("Validation failed: %s", err)
		return errors.New("10003")
	}

	if disbursement.LoanID <= 0 {
		log.Printf("LoanID must be provided and greater than zero")
		return errors.New("10003")
	}

	if !disbursement.DisbursementStatus.IsValid() {
		log.Printf("Invalid DisbursementStatus: %v", disbursement.DisbursementStatus)
		return errors.New("10003")
	}

	if disbursement.StaffID <= 0 {
		log.Printf("StaffID must be provided and greater than zero")
		return errors.New("10003")
	}

	if disbursement.SignedAgreementURL == "" {
		log.Printf("SignedAgreementURL must be provided")
		return errors.New("10003")
	}

	log.Printf("Validation passed for update loan disbursement request: LoanID=%d", disbursement.LoanID)
	return nil
}

func (l LoanDisbursementValidatorImpl) ValidateTransitionStatus(from interface{}, to interface{}) bool {

	var existing enum.LoanDisbursementStatus
	var newStatus enum.LoanDisbursementStatus
	existing = from.(enum.LoanDisbursementStatus)
	newStatus = to.(enum.LoanDisbursementStatus)

	log.Printf("Validating status transition: ExistingStatus=%v, NewStatus=%v", existing, newStatus)
	validTransitions := map[enum.LoanDisbursementStatus][]enum.LoanDisbursementStatus{
		enum.LoanDisbursementPending:   {enum.LoanDisbursementCompleted, enum.LoanDisbursementCancelled},
		enum.LoanDisbursementCompleted: {}, // No transitions from Success
		enum.LoanDisbursementCancelled: {}, // No transitions from Failed
	}

	allowedStatuses, ok := validTransitions[existing]
	if !ok {
		log.Printf("No valid transitions found for status: %v", existing)
		return false
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			log.Printf("Valid transition found: %v -> %v", existing, newStatus)
			return true
		}
	}
	log.Printf("Invalid transition: %v -> %v", existing, newStatus)
	return false
}
