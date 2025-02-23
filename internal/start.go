package internal

import (
	"github.com/test/loan-service/internal/handler/api"
	"github.com/test/loan-service/internal/handler/kafka"
	"github.com/test/loan-service/internal/handler/middleware"
	"net/http"

	"github.com/labstack/echo"
	"github.com/test/loan-service/internal/infra"
	"go.uber.org/dig"
)

func Start(
	di *dig.Container,
	cfg *infra.EchoCfg,
	e *echo.Echo,
) (err error) {

	e.Use(middleware.I18nMiddleware)
	e.Use(middleware.ErrorHandlerMiddleware)
	e.Use(middleware.SuccessHandlerMiddleware)

	if err = di.Invoke(api.NewLoanHandler); err != nil {
		return err
	}
	if err = di.Invoke(api.NewLoanApprovalHandler); err != nil {
		return err
	}

	if err = di.Invoke(api.NewLoanFundingHandler); err != nil {
		return err
	}
	if err = di.Invoke(api.NewLoanDisbursementHandler); err != nil {
		return err
	}

	if err = di.Invoke(kafka.NewKafkaHandler); err != nil {
		return err
	}

	return e.StartServer(&http.Server{
		Addr:         cfg.Address,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})
}
