package infra

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	repo "github.com/test/loan-service/internal/repository"
	"github.com/test/loan-service/internal/service"
	"github.com/typical-go/typical-go/pkg/typapp"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal(err.Error())
	}

	// config properties
	typapp.Provide("", LoadDatabaseCfg)
	typapp.Provide("", LoadKafkaCfg)
	typapp.Provide("", LoadEchoCfg)
	typapp.Provide("", LoadSMTPConfig)

	// config
	typapp.Provide("", NewDatabases)
	typapp.Provide("", NewKafkaClients)
	typapp.Provide("", NewEcho)
	typapp.Provide("", NewSMTPs)

	// repo dependency injection
	typapp.Provide("", repo.NewLoanRepo)
	typapp.Provide("", repo.NewLoanDetailRepo)
	typapp.Provide("", repo.NewLoanApprovalRepo)
	typapp.Provide("", repo.NewApprovalDocumentRepo)
	typapp.Provide("", repo.NewLoanFundingRepo)
	typapp.Provide("", repo.NewLoanDisbursementRepo)

	// service dependency injection
	typapp.Provide("", service.NewLoanSvc)
	typapp.Provide("", service.NewEmailSvc)
	typapp.Provide("", service.NewLoanDisbursementSvc)
	typapp.Provide("", service.NewLoanApprovalSvc)
	typapp.Provide("", service.NewLoanDetailSvc)
	typapp.Provide("", service.NewLoanFundingSvc)

}
