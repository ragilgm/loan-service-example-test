package models

import (
	"github.com/test/loan-service/internal/dto"
	"github.com/test/loan-service/internal/enum"
)

type (
	LoanDetailRequest struct {
		Detail     *dto.LoanDetailRequestDTO
		BorrowerID int64
		LoanID     int64
	}

	LoanApprovalRequest struct {
		Page   uint64
		Size   uint64
		Status *enum.ApprovalStatus
	}

	LoanDisbursementRequest struct {
		Page   uint64
		Size   uint64
		Status *enum.LoanDisbursementStatus
	}

	LoanRequest struct {
		Page   uint64
		Size   uint64
		Status *enum.LoanStatus
	}
)
