package enum

type LoanFundingStatus string

const (
	LoanFundingPending   LoanFundingStatus = "pending"
	LoanFundingInvested  LoanFundingStatus = "invested"
	LoanFundingFailed    LoanFundingStatus = "failed"
	LoanFundingOngoing   LoanFundingStatus = "on_going"
	LoanFundingCompleted LoanFundingStatus = "completed"
)

func (s LoanFundingStatus) IsValid() bool {
	switch s {
	case LoanFundingPending, LoanFundingInvested, LoanFundingFailed, LoanFundingOngoing, LoanFundingCompleted:
		return true
	}
	return false
}
