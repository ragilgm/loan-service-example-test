package message

import "github.com/test/loan-service/internal/enum"

type UpdateLoanMessage struct {
	LoanID     int64
	LoanStatus enum.LoanStatus
}
