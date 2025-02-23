package enum

type LoanDisbursementStatus string

const (
	LoanDisbursementPending   LoanDisbursementStatus = "pending"
	LoanDisbursementCompleted LoanDisbursementStatus = "completed"
	LoanDisbursementCancelled LoanDisbursementStatus = "cancelled"
)

// IsValid checks if the LoanDisbursementStatus is valid.
func (s LoanDisbursementStatus) IsValid() bool {
	switch s {
	case LoanDisbursementPending, LoanDisbursementCompleted, LoanDisbursementCancelled:
		return true
	}
	return false
}
