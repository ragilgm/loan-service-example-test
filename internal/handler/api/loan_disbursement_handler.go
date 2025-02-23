package api

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/test/loan-service/internal/dto"
	"github.com/test/loan-service/internal/enum"
	"github.com/test/loan-service/internal/handler/middleware"
	"github.com/test/loan-service/internal/service"
	"go.uber.org/dig"
	"net/http"
	"strconv"
)

type (
	LoanDisbursementHandler struct {
		dig.In
		loanDisbursementSvc service.LoanDisbursementSvc
	}
)

func NewLoanDisbursementHandler(e *echo.Echo, loanDisbursementSvc service.LoanDisbursementSvc) *LoanDisbursementHandler {
	handler := &LoanDisbursementHandler{
		loanDisbursementSvc: loanDisbursementSvc,
	}

	// Define the routes
	e.GET("/loan-disbursements/:id", handler.GetByID)
	e.GET("/loan-disbursements", handler.GetAll)
	e.PUT("/loan-disbursements/:id", handler.Update)

	return handler
}

// Create - Handler for creating loan disbursement
func (ldh *LoanDisbursementHandler) Create(c echo.Context) error {
	var request dto.LoanDisbursementRequestDTO
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()

	// Call the service to create loan disbursement
	err = ldh.loanDisbursementSvc.Create(ctx, &request)
	if err != nil {
		return err
	}
	return middleware.SendSuccess(c, nil)
}

// GetByID - Handler to get loan disbursement by ID
func (ldh *LoanDisbursementHandler) GetByID(c echo.Context) error {
	disbursementIDStr := c.Param("id")
	disbursementID, err := strconv.ParseInt(disbursementIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseError{Message: "Invalid disbursement ID"})
	}

	ctx := c.Request().Context()

	// Call the service to get loan disbursement by ID
	loanDisbursement, err := ldh.loanDisbursementSvc.GetByID(ctx, disbursementID)
	if err != nil {
		return err
	}

	return middleware.SendSuccess(c, loanDisbursement)
}

// GetAll - Handler to get all loan disbursements with pagination and optional status filter
func (ldh *LoanDisbursementHandler) GetAll(c echo.Context) error {
	pageStr := c.QueryParam("page")
	sizeStr := c.QueryParam("size")
	statusStr := c.QueryParam("disbursement_status")

	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil || page == 0 {
		page = 1
	}

	size, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil || size == 0 {
		size = 10
	}

	var disburseStatus enum.LoanDisbursementStatus

	// Validasi apakah status yang diberikan valid
	if statusStr != "" {
		disburseStatus = enum.LoanDisbursementStatus(statusStr)

		// Validasi status sesuai dengan enum
		if !disburseStatus.IsValid() {
			return errors.New("10002")
		}
	}

	// Prepare request object for the service
	request := service.LoanDisbursementRequest{
		Page:   page,
		Size:   size,
		Status: &disburseStatus,
	}

	ctx := c.Request().Context()

	// Call the service to get loan disbursements with pagination
	loanDisbursements, totalRecords, err := ldh.loanDisbursementSvc.GetAllPage(ctx, request)
	if err != nil {
		return err
	}

	return middleware.SendSuccess(c, dto.PaginationHelper(loanDisbursements, totalRecords, int(page), int(size)))

}

func (ldh *LoanDisbursementHandler) Update(c echo.Context) error {
	// Parse disbursement ID from URL
	disbursementIDStr := c.Param("id")
	disbursementID, err := strconv.ParseInt(disbursementIDStr, 10, 64)
	if err != nil {
		return middleware.SendSuccess(c, "Error updating LoanDisbursement")
	}

	var request dto.UpdateLoanDisbursementRequestDTO
	// Bind the request body to the DTO
	err = c.Bind(&request)
	if err != nil {
		return err
	}
	// Call the service to update the loan disbursement
	ctx := c.Request().Context()
	err = ldh.loanDisbursementSvc.Update(ctx, disbursementID, &request)
	if err != nil {
		return err
	}

	// Return a successful response
	return middleware.SendSuccess(c, "Loan disbursement updated")
}
