package dto

import "time"

type LoanFundingRequestDTO struct {
	OrderNumber        string  `json:"order_number" validate:"required"`
	LoanID             int64   `json:"loan_id," validate:"required"`
	LenderID           int64   `json:"lender_id" validate:"required"`
	LenderEmail        string  `json:"lender_email" validate:"required,email"`
	InvestmentAmount   float64 `json:"investment_amount" validate:"required"`
	LenderAgreementURL string  `json:"lender_agreement_url" validate:"required"`
}

type LoanFundingResponseDTO struct {
	ID                 int64      `json:"id"`
	LoanOrderNumber    string     `json:"loan_order_number"`
	OrderNumber        string     `json:"order_number"`
	LoanID             int64      `json:"loan_id"`
	LenderID           int64      `json:"lender_id"`
	LenderEmail        string     `json:"lender_email"`
	InvestmentAmount   float64    `json:"investment_amount"`
	Rate               float64    `json:"rate"`
	Interest           float64    `json:"interest"`
	ROI                float64    `json:"roi"`
	InterestPaid       float64    `json:"interest_paid"`
	CapitalAmountPaid  float64    `json:"capital_amount_paid"`
	TotalAmountPaid    float64    `json:"total_amount_paid"`
	InvestmentDate     time.Time  `json:"investment_date"`
	Status             string     `json:"status"`
	LenderAgreementURL string     `json:"lender_agreement_url"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}
