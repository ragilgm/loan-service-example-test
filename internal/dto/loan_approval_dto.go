package dto

import (
	"github.com/test/loan-service/internal/enum"
	"time"
)

type LoanApprovalRequestDTO struct {
	LoanID int64 `json:"loan_id" valid:"required"` // Loan ID
}

type UpdateLoanApprovalRequestDTO struct {
	StaffID           int64                        `json:"staff_id" valid:"required"`
	ApprovalStatus    enum.ApprovalStatus          `json:"approval_status" valid:"required"`
	ApprovalDocuments []ApprovalDocumentRequestDTO `json:"approval_documents" valid:"required"`
}

type LoanApprovalResponseDTO struct {
	ID               int64                          `json:"id"`
	LoanID           int64                          `json:"loan_id"`
	ApprovalNumber   string                         `json:"approval_number"`
	StaffID          *int64                         `json:"staff_id,omitempty"`
	ApprovalDate     *time.Time                     `json:"approval_date,omitempty"`
	ApprovalStatus   enum.ApprovalStatus            `json:"approval_status"`
	CreatedAt        time.Time                      `json:"created_at"`
	UpdatedAt        time.Time                      `json:"updated_at"`
	DeletedAt        *time.Time                     `json:"deleted_at,omitempty"`
	ApprovalDocument *[]ApprovalDocumentResponseDTO `json:"approval_document,omitempty"`
}
