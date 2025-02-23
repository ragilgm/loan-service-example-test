package enum

type LoanStatus string

const (
	Proposed  LoanStatus = "proposed"
	Rejected  LoanStatus = "rejected"
	Approved  LoanStatus = "approved"
	Invested  LoanStatus = "invested"
	Disbursed LoanStatus = "disbursed"
	Completed LoanStatus = "completed"
)

func (s LoanStatus) IsValid() bool {
	switch s {
	case Proposed, Rejected, Approved, Invested, Disbursed, Completed:
		return true
	}
	return false
}

type LoanType string

const (
	Productive  LoanType = "productive"
	Consumptive LoanType = "consumptive"
)

func (s LoanType) IsValid() bool {
	switch s {
	case Productive, Consumptive:
		return true
	}
	return false
}
