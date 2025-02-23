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
	LoanCtrlImpl struct {
		dig.In
		loanSvc service.LoanSvc
	}
)

func NewLoanHandler(e *echo.Echo, loanSvc service.LoanSvc) *LoanCtrlImpl {
	handler := &LoanCtrlImpl{
		loanSvc: loanSvc,
	}

	e.POST("/loans", handler.Create)
	e.GET("/loans", handler.GetAll)
	e.GET("/loans/:id", handler.GetByID)

	return handler
}

func (ic LoanCtrlImpl) Create(c echo.Context) (err error) {
	var loanRequest dto.LoanRequestDTO
	err = c.Bind(&loanRequest)
	if err != nil {
		return errors.New("10002")
	}
	ctx := c.Request().Context()

	_, err = ic.loanSvc.Create(ctx, &loanRequest)
	if err != nil {
		return err
	}

	return middleware.SendSuccess(c, "Loan created")
}

func (ic *LoanCtrlImpl) GetAll(c echo.Context) (err error) {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1 // Default to page 1
	}

	size, err := strconv.Atoi(c.QueryParam("size"))
	if err != nil || size <= 0 {
		size = 10 // Default to 10 items per page
	}

	// Ambil nilai parameter status dari query string
	status := c.QueryParam("loan_status")
	var loanStatus enum.LoanStatus

	// Validasi apakah status yang diberikan valid
	if status != "" {
		loanStatus = enum.LoanStatus(status)

		// Validasi status sesuai dengan enum
		if !loanStatus.IsValid() {
			return errors.New("10002")
		}
	}

	ctx := c.Request().Context()

	request := service.LoanRequest{
		Page:   uint64(page),
		Size:   uint64(size),
		Status: &loanStatus,
	}

	// Get paginated loans
	loans, totalRecords, err := ic.loanSvc.GetAllPage(ctx, request)
	if err != nil {
		return err
	}

	// Return paginated response
	return middleware.SendSuccess(c, dto.PaginationHelper(loans, totalRecords, page, size))
}

func (ic LoanCtrlImpl) GetByID(c echo.Context) (err error) {
	// Get the loan ID from the URL parameter
	loanIDStr := c.Param("id")
	loanID, err := strconv.ParseInt(loanIDStr, 10, 64)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()

	// Call service to get loan by ID
	loan, err := ic.loanSvc.GetByID(ctx, loanID)
	if err != nil {
		return err
	}

	// Return the loan details
	return middleware.SendSuccess(c, loan)
}
