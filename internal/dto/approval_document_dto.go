package dto

import "time"

type ApprovalDocumentRequestDTO struct {
	DocumentType string  `json:"document_type" validate:"required"`
	FileURL      string  `json:"file_url" validate:"required"`
	Description  *string `json:"description,omitempty"`
}

type ApprovalDocumentResponseDTO struct {
	ID             int64      `json:"id"`
	LoanApprovalID int64      `json:"loan_approval_id"`
	DocumentType   string     `json:"document_type"`
	FileURL        string     `json:"file_url"`
	Description    *string    `json:"description,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
