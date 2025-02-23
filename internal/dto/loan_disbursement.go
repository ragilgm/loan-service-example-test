package dto

import (
	"github.com/test/loan-service/internal/enum"
	"time"
)

type LoanDisbursementRequestDTO struct {
	LoanID         int64   `json:"loan_id" validate:"required"`
	DisburseAmount float64 `json:"disburse_amount" validate:"required,gt=0"`
}

type UpdateLoanDisbursementRequestDTO struct {
	LoanID             int64                       `json:"loan_id" validate:"required"`
	DisbursementStatus enum.LoanDisbursementStatus `json:"disbursement_status" validate:"required"`  // Status (Pending, Completed, etc.)
	StaffID            int64                       `json:"staff_id" validate:"required"`             // Staff ID handling the disbursement
	SignedAgreementURL string                      `json:"signed_agreement_url" validate:"required"` // URL to signed agreement

}

type LoanDisbursementResponseDTO struct {
	ID                 int64                       `json:"id"`                             // Disbursement ID
	LoanID             int64                       `json:"loan_id"`                        // Loan ID
	DisburseCode       string                      `json:"disburse_code"`                  // Disbursement code
	DisburseAmount     float64                     `json:"disburse_amount"`                // Disbursed amount
	DisbursementStatus enum.LoanDisbursementStatus `json:"disbursement_status"`            // Status (Pending, Completed, etc.)
	DisburseDate       *time.Time                  `json:"disburse_date,omitempty"`        // Disbursement date
	StaffID            *int64                      `json:"staff_id,omitempty"`             // Staff ID handling the disbursement
	AgreementURL       string                      `json:"agreement_url"`                  // URL template agreement url
	SignedAgreementURL *string                     `json:"signed_agreement_url,omitempty"` // URL to signed agreement
	CreatedAt          time.Time                   `json:"created_at"`                     // Date of creation
	UpdatedAt          time.Time                   `json:"updated_at"`                     // Date of last update
	DeletedAt          *time.Time                  `json:"deleted_at,omitempty"`           // Date of deletion if applicable
}
