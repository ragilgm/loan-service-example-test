package api

import (
	"errors"
	"github.com/test/loan-service/internal/enum"
	"github.com/test/loan-service/internal/handler/middleware"
	"strconv"

	"github.com/labstack/echo"
	"github.com/test/loan-service/internal/dto"
	"github.com/test/loan-service/internal/service"
	"go.uber.org/dig"
)

type (
	LoanApprovalHandler struct {
		dig.In
		approvalSvc service.LoanApprovalSvc
	}
)

func NewLoanApprovalHandler(e *echo.Echo, approvalSvc service.LoanApprovalSvc) *LoanApprovalHandler {
	handler := &LoanApprovalHandler{
		approvalSvc: approvalSvc,
	}

	e.GET("/loans/approvals", handler.GetAllPage)
	e.PUT("/loans/approvals/:id", handler.Update)

	return handler
}

func (ic *LoanApprovalHandler) GetAllPage(c echo.Context) (err error) {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1 // Default to page 1
	}

	size, err := strconv.Atoi(c.QueryParam("size"))
	if err != nil || size <= 0 {
		size = 10 // Default to 10 items per page
	}

	// Ambil nilai parameter status dari query string
	status := c.QueryParam("approval_status")
	var loanStatus enum.ApprovalStatus

	// Validasi apakah status yang diberikan valid
	if status != "" {
		loanStatus = enum.ApprovalStatus(status)

		// Validasi status sesuai dengan enum
		if !loanStatus.IsValid() {
			return errors.New("10002")
		}
	}

	ctx := c.Request().Context()

	request := service.LoanApprovalRequest{
		Page:   uint64(page),
		Size:   uint64(size),
		Status: &loanStatus,
	}

	// Get paginated loan approval
	loans, totalRecords, err := ic.approvalSvc.GetAllPage(ctx, request)
	if err != nil {
		return err
	}

	// Return paginated response
	return middleware.SendSuccess(c, dto.PaginationHelper(loans, totalRecords, page, size))
}

func (ic *LoanApprovalHandler) Update(c echo.Context) error {

	approvalIdStr := c.Param("id")
	approvalID, err := strconv.ParseInt(approvalIdStr, 10, 64)
	if err != nil {
		return errors.New("10002")
	}

	var loanRequest dto.UpdateLoanApprovalRequestDTO
	err = c.Bind(&loanRequest)
	if err != nil {
		return errors.New("10002")
	}

	ctx := c.Request().Context()

	err = ic.approvalSvc.Update(ctx, approvalID, &loanRequest)
	if err != nil {
		return err
	}

	return middleware.SendSuccess(c, "Loan approval updated")
}
