package service

import (
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

const (
	SMS_MAX_LENGTH int = 160
	SMS_MIN_LENGTH int = 1
)

var (
	ErrSMSGeneric        = internal.NewError("SMS delivery failed")
	ErrSMSInvalidMessage = internal.NewError("SMS message invalid")
)

type SMS struct {
	class.Service
	client twilio.RestClient
}

func NewSMS(configuration internal.Configuration, logger core.Logger) *SMS {
	client := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username:   configuration.TwilioApiKey,
		Password:   configuration.TwilioApiSecret,
		AccountSid: configuration.TwilioAccountSID,
	})

	client.SetEdge(configuration.TwilioEdge)
	client.SetTimeout(time.Duration(configuration.GracefulTimeout) * time.Second)

	return &SMS{
		Service: *class.NewService(configuration, logger),
		client:  *client,
	}
}

func (self *SMS) Send(user model.User, message string) error {
	if len(message) < SMS_MIN_LENGTH || len(message) > SMS_MAX_LENGTH {
		return ErrSMSInvalidMessage()
	}

	if self.Configuration.Environment == internal.Environment.PRODUCTION {
		return self.sendReal(user, message)
	} else {
		return self.sendFake(user, message)
	}
}

func (self *SMS) sendFake(user model.User, message string) error {
	self.Logger.Debugf("SMS sent to: %s: %s", user, message)
	return nil
}

func (self *SMS) sendReal(user model.User, message string) error {
	params := &openapi.CreateMessageParams{}
	params.SetFrom(self.Configuration.TwilioFromPhone)
	params.SetTo(user.Phone)
	params.SetBody(message)

	_, err := self.client.ApiV2010.CreateMessage(params)
	if err != nil {
		return ErrSMSGeneric().Wrap(err)
	}

	return nil
}
