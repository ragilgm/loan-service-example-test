package enum

type ApprovalStatus string

const (
	ApprovalPending  ApprovalStatus = "pending"
	ApprovalApproved ApprovalStatus = "approved"
	ApprovalRejected ApprovalStatus = "rejected"
)

func (s ApprovalStatus) IsValid() bool {
	switch s {
	case ApprovalPending, ApprovalApproved, ApprovalRejected:
		return true
	}
	return false
}
