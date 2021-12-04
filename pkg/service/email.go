package service

import (
	"context"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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
	client sendgrid.Client
}

func NewEmailService(configuration internal.Configuration, logger core.Logger) *EmailService {
	return &EmailService{
		Service: *class.NewService(configuration, logger),
		client:  *sendgrid.NewSendClient(configuration.SendGridApiKey),
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
		go func() {
			self.Logger.Error(self.sendReal(receiverEmail, subject, body))
		}()
		return nil
	} else {
		return self.sendFake(receiverEmail, subject, body)
	}
}

func (self *EmailService) sendFake(receiverEmail string, subject string, body string) error {
	self.Logger.Debugf("Email sent to: %s: %s", receiverEmail, subject)
	return nil
}

func (self *EmailService) sendReal(receiverEmail string, subject string, body string) error {
	from := mail.NewEmail(self.Configuration.SendGridFromName, self.Configuration.SendGridFromEmail)
	to := mail.NewEmail(strings.Title(receiverEmail[:strings.Index(receiverEmail, "@")]), receiverEmail)

	message := mail.NewSingleEmail(from, subject, to, body, body)

	_, err := self.client.Send(message)
	if err != nil {
		return ErrEmailGeneric().Wrap(err)
	}

	return nil
}

func (self *EmailService) Close(ctx context.Context) error {
	self.Logger.Info("Closing email service")

	return nil
}
