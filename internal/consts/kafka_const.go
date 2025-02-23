package consts

type KafkaTopic string

const (
	ApprovalLoanTopic   KafkaTopic = "loan-approval-topic"
	LoanDisburseTopic   KafkaTopic = "loan-disburse-topic"
	FundingProcessTopic KafkaTopic = "funding-process-topic"
)
