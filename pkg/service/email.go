package service

import (
	"context"

	"github.com/badoux/checkmail"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
)

const (
	EMAIL_SUBJECT_MAX_LENGTH = 70
	EMAIL_SUBJECT_MIN_LENGTH = 1
)

var (
	ErrEmailGeneric        = internal.NewError("Email delivery failed")
	ErrEmailInvalid        = internal.NewError("Email receiver invalid")
	ErrEmailInvalidSubject = internal.NewError("Email subject invalid")
)

type EmailService struct {
	class.Service
}

func NewEmailService(configuration internal.Configuration, logger core.Logger) *EmailService {
	return &EmailService{
		Service: *class.NewService(configuration, logger),
	}
}

func (self *EmailService) Send(receiverEmail string, subject string, body string) error {
	if len(subject) < EMAIL_SUBJECT_MIN_LENGTH || len(subject) > EMAIL_SUBJECT_MAX_LENGTH {
		return ErrEmailInvalidSubject()
	}

	err := checkmail.ValidateFormat(receiverEmail)
	if err != nil {
		return ErrEmailInvalid().Wrap(err)
	}

	if self.Configuration.Environment == internal.Environment.PRODUCTION {
		return self.sendReal(receiverEmail, subject, body)
	} else {
		return self.sendFake(receiverEmail, subject, body)
	}
}

func (self *EmailService) sendFake(receiverEmail string, subject string, body string) error {
	self.Logger.Debugf("Email sent to: %s: %s", receiverEmail, subject)
	return nil
}

// https://docs.sendgrid.com/for-developers/sending-email/quickstart-go
func (self *EmailService) sendReal(receiverEmail string, subject string, body string) error {
	return nil
}

func (self *EmailService) Close(ctx context.Context) error {
	self.Logger.Info("Closing email service")

	return nil
}
