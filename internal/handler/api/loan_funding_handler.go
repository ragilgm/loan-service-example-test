package api

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/test/loan-service/internal/dto"
	"github.com/test/loan-service/internal/service"
	"go.uber.org/dig"
	"net/http"
	"strconv"
)

type (
	LoanFundingHandler struct {
		dig.In
		loanFundingSvc service.LoanFundingSvc
	}
)

func NewLoanFundingHandler(e *echo.Echo, loanFundingSvc service.LoanFundingSvc) *LoanFundingHandler {
	handler := &LoanFundingHandler{
		loanFundingSvc: loanFundingSvc,
	}

	e.POST("/loan-fundings", handler.Create)
	e.GET("/loan-fundings/:id", handler.GetByID)
	e.GET("/loan-fundings/lender/:lender_id", handler.GetByLenderID)

	return handler
}

// Create - Handler for creating loan funding
func (lh *LoanFundingHandler) Create(c echo.Context) error {
	var request dto.LoanFundingRequestDTO
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()

	// Call the service to create loan funding
	err = lh.loanFundingSvc.Create(ctx, &request)
	if err != nil {
		return err
	}

	return dto.SendSuccess(c, "Loan funding created")
}

// GetByID - Handler to get loan funding by ID
func (lh *LoanFundingHandler) GetByID(c echo.Context) error {
	loanFundingIDStr := c.Param("id")
	loanFundingID, err := strconv.ParseInt(loanFundingIDStr, 10, 64)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()

	// Call the service to get loan funding by ID
	loanFunding, err := lh.loanFundingSvc.GetByID(ctx, loanFundingID)
	if err != nil {
		return err
	}

	return dto.SendSuccess(c, loanFunding)
}

// GetByLenderID - Handler to get all loan fundings by lender ID
func (lh *LoanFundingHandler) GetByLenderID(c echo.Context) error {
	lenderIDStr := c.Param("lender_id")
	lenderID, err := strconv.ParseInt(lenderIDStr, 10, 64)
	if err != nil {
		return errors.New("10002")
	}

	ctx := c.Request().Context()

	// Call the service to get loan fundings by lender ID
	loanFundings, err := lh.loanFundingSvc.GetByLenderID(ctx, lenderID)
	if err != nil {
		return err
	}

	return dto.SendSuccess(c, loanFundings)
}
