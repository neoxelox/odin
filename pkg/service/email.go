package service

import (
	"context"

	"github.com/badoux/checkmail"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
)

var (
	ErrEmailGeneric = internal.NewError("Email delivery failed")
	ErrEmailInvalid = internal.NewError("Email receiver invalid")
)

type EmailService struct {
	class.Service
}

func NewEmailService(configuration internal.Configuration, logger core.Logger) *EmailService {
	return &EmailService{
		Service: *class.NewService(configuration, logger),
	}
}

func (self *EmailService) Send(receiverEmail string, message string) error {
	err := checkmail.ValidateFormat(receiverEmail)
	if err != nil {
		return ErrEmailInvalid().Wrap(err)
	}

	err = checkmail.ValidateHost(receiverEmail)
	if err != nil {
		return ErrEmailInvalid().Wrap(err)
	}

	if self.Configuration.Environment == internal.Environment.PRODUCTION {
		return self.sendReal(receiverEmail, message)
	} else {
		return self.sendFake(receiverEmail, message)
	}
}

func (self *EmailService) sendFake(receiverEmail string, message string) error {
	self.Logger.Debugf("Email sent to: %s: %s", receiverEmail, message)
	return nil
}

// https://docs.sendgrid.com/for-developers/sending-email/quickstart-go
func (self *EmailService) sendReal(receiverEmail string, message string) error {
	return nil
}

func (self *EmailService) Close(ctx context.Context) error {
	self.Logger.Info("Closing email service")

	return nil
}
