package infra

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	repo "github.com/test/loan-service/internal/repository"
	"github.com/test/loan-service/internal/service"
	"github.com/test/loan-service/internal/service/validator"
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

	// validator dependency injection
	typapp.Provide("loan_validator", validator.NewLoanValidator)
	typapp.Provide("loan_approval_validator", validator.NewLoanApprovalValidator)
	typapp.Provide("loan_funding_validator", validator.NewLoanFundingValidator)
	typapp.Provide("loan_detail_validator", validator.NewLoanDetailValidator)
	typapp.Provide("loan_disbursement_validator", validator.NewLoanDisbursementValidator)

	// service dependency injection
	typapp.Provide("", service.NewLoanSvc)
	typapp.Provide("", service.NewEmailSvc)
	typapp.Provide("", service.NewLoanDisbursementSvc)
	typapp.Provide("", service.NewLoanApprovalSvc)
	typapp.Provide("", service.NewLoanDetailSvc)
	typapp.Provide("", service.NewLoanFundingSvc)

}
