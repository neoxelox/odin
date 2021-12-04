package service

import (
	"context"
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/nyaruka/phonenumbers"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

const (
	SMS_MAX_LENGTH = 160
	SMS_MIN_LENGTH = 1
)

var (
	ErrSMSGeneric        = internal.NewError("SMS delivery failed")
	ErrSMSInvalidPhone   = internal.NewError("SMS receiver phone invalid")
	ErrSMSInvalidMessage = internal.NewError("SMS message invalid")
)

type SMSService struct {
	class.Service
	client twilio.RestClient
}

func NewSMSService(configuration internal.Configuration, logger core.Logger) *SMSService {
	client := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username:   configuration.TwilioApiKey,
		Password:   configuration.TwilioApiSecret,
		AccountSid: configuration.TwilioAccountSID,
	})

	client.SetEdge(configuration.TwilioEdge)
	client.SetTimeout(time.Duration(configuration.GracefulTimeout) * time.Second)

	return &SMSService{
		Service: *class.NewService(configuration, logger),
		client:  *client,
	}
}

func (self *SMSService) Send(receiverPhone string, message string) error {
	if len(message) < SMS_MIN_LENGTH || len(message) > SMS_MAX_LENGTH {
		return ErrSMSInvalidMessage()
	}

	ph, err := phonenumbers.Parse(receiverPhone, "ES")
	if err != nil {
		return ErrSMSInvalidPhone().Wrap(err)
	}

	if !phonenumbers.IsValidNumber(ph) {
		return ErrSMSInvalidPhone()
	}

	if self.Configuration.Environment == internal.Environment.PRODUCTION && self.Configuration.ServiceSMSEnabled {
		go func() {
			self.Logger.Error(self.sendReal(receiverPhone, message))
		}()
		return nil
	} else {
		return self.sendFake(receiverPhone, message)
	}
}

func (self *SMSService) sendFake(receiverPhone string, message string) error {
	self.Logger.Debugf("SMS sent to: %s: %s", receiverPhone, message)
	return nil
}

func (self *SMSService) sendReal(receiverPhone string, message string) error {
	params := &openapi.CreateMessageParams{}
	params.SetFrom(self.Configuration.TwilioFromPhone)
	params.SetTo(receiverPhone)
	params.SetBody(message)

	_, err := self.client.ApiV2010.CreateMessage(params)
	if err != nil {
		return ErrSMSGeneric().Wrap(err)
	}

	return nil
}

func (self *SMSService) Close(ctx context.Context) error {
	self.Logger.Info("Closing SMS service")

	return nil
}
