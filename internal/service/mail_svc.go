package service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"gopkg.in/gomail.v2"
)

type (
	SendEmailInput struct {
		To      []string
		Subject string
		Body    string
	}

	EmailSvc interface {
		SendEmail(ctx context.Context, input SendEmailInput) error
	}

	EmailSvcImpl struct {
		dig.In
		Mailer *gomail.Dialer
	}
)

func NewEmailSvc(impl EmailSvcImpl) EmailSvc {
	return &impl
}

// SendEmail untuk mengirim email
func (s *EmailSvcImpl) SendEmail(ctx context.Context, input SendEmailInput) error {
	// Membuat pesan email baru
	message := gomail.NewMessage()
	message.SetHeader("From", "ragilnamasaya@gmail.com") // Gantilah dengan email pengirim
	message.SetHeader("To", input.To...)
	message.SetHeader("Subject", input.Subject)
	message.SetBody("text/plain", input.Body)

	// Mencoba untuk mengirim email
	if err := s.Mailer.DialAndSend(message); err != nil {
		logrus.Errorf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logrus.Infof("Email sent successfully to %v", input.To)
	return nil
}
