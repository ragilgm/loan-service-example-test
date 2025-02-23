package infra

import (
	"fmt"
	"go.uber.org/dig"
	"gopkg.in/gomail.v2"
	"strconv"
)

type (
	SMTPs struct {
		dig.Out
		Mailer *gomail.Dialer
	}

	SMTPCfgs struct {
		dig.In
		SMTP *SMTPCfg
	}

	SMTPCfg struct {
		Host     string `envconfig:"SMTP_HOST" required:"true" default:"smtp.gmail.com"`
		Port     string `envconfig:"SMTP_PORT" required:"true" default:"587"`
		Username string `envconfig:"SMTP_USERNAME" required:"true"`
		Password string `envconfig:"SMTP_PASSWORD" required:"true"`
	}
)

// NewSMTPs creates a new instance of SMTP (using gomail Dialer)
func NewSMTPs(cfgs SMTPCfgs) SMTPs {
	return SMTPs{
		Mailer: openSMTP(cfgs.SMTP),
	}
}

// openSMTP initializes a gomail Dialer to connect to the SMTP server
func openSMTP(cfg *SMTPCfg) *gomail.Dialer {
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		panic(err)
	}
	// Set up the SMTP server configuration
	dialer := gomail.NewDialer(cfg.Host, port, cfg.Username, cfg.Password)
	dialer.SSL = false // For TLS, set to true if required by the mail provider

	// Return the created dialer (SMTP client)
	return dialer
}

func (s *SMTPs) CheckSMTPConnection() error {
	// Try to dial to the SMTP server to check the connection
	_, err := s.Mailer.Dial()
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}

	return nil
}
