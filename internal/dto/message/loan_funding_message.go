package message

type FundingProcessMessage struct {
	LoanOrderNumber string `json:"loan_order_number"`
	LoanID          int64  `json:"loan_id"`
}
