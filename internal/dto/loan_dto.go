package dto

import (
	"github.com/test/loan-service/internal/enum"
	"time"
)

type LoanRequestDTO struct {
	BorrowerID    int64                `json:"borrower_id" valid:"required"`
	RequestAmount float64              `json:"request_amount" valid:"required"`
	LoanGrade     string               `json:"loan_grade" valid:"required"`
	LoanType      enum.LoanType        `json:"loan_type" valid:"required"`
	Rate          float64              `json:"rate" valid:"required"`
	Tenures       int                  `json:"tenures" valid:"required"`
	Detail        LoanDetailRequestDTO `json:"detail" valid:"required"`
}

type LoanResponseDTO struct {
	ID                   int64                  `json:"id"`                         // Loan ID
	LoanCode             string                 `json:"loan_code"`                  // Loan code
	BorrowerID           int64                  `json:"borrower_id"`                // Borrower ID
	RequestAmount        float64                `json:"request_amount"`             // Loan request amount
	LoanGrade            string                 `json:"loan_grade"`                 // Loan grade (A, B, C, D)
	LoanType             enum.LoanType          `json:"loan_type"`                  // Type of loan (productive, consumptive, etc.)
	TotalInvestedAmount  float64                `json:"total_invested_amount"`      // Total amount invested
	InvestorCount        int64                  `json:"investor_count"`             // Number of investors participating
	FundingDeadline      *time.Time             `json:"funding_deadline,omitempty"` // Funding deadline
	LoanStatus           enum.LoanStatus        `json:"loan_status"`                // Loan status (proposed, rejected, approved, invested)
	Rate                 float64                `json:"rate"`                       // Interest rate
	Tenures              int64                  `json:"tenures"`                    // Loan tenure
	TotalRepaymentAmount float64                `json:"total_repayment_amount"`     // Total repayment amount needed
	InvestmentPercentage float64                `json:"investment_percentage"`      // Investor profit sharing percentage
	AgreementLetterLink  string                 `json:"agreement_letter_link"`      // Link to generated loan agreement letter
	CreatedAt            time.Time              `json:"created_at"`                 // Loan creation date
	UpdatedAt            time.Time              `json:"updated_at"`                 // Loan status update date
	DeletedAt            *time.Time             `json:"deleted_at,omitempty"`       // Loan deletion date (if applicable)
	LoanDetail           *LoanDetailResponseDTO `json:"loan_detail,omitempty"`
}
