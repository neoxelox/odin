package service

import (
	"context"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
)

const (
	PUSH_MAX_LENGTH = 240
	PUSH_MIN_LENGTH = 1
)

var (
	ErrPushGeneric        = internal.NewError("Push delivery failed")
	ErrPushInvalidMessage = internal.NewError("Push message invalid")
)

type PushService struct {
	class.Service
}

func NewPushService(configuration internal.Configuration, logger core.Logger) *PushService {
	return &PushService{
		Service: *class.NewService(configuration, logger),
	}
}

func (self *PushService) Send(receiverSubscription string, message string) error {
	if len(message) < SMS_MIN_LENGTH || len(message) > SMS_MAX_LENGTH {
		return ErrPushInvalidMessage()
	}

	if self.Configuration.Environment == internal.Environment.PRODUCTION {
		return self.sendReal(receiverSubscription, message)
	} else {
		return self.sendFake(receiverSubscription, message)
	}
}

func (self *PushService) sendFake(receiverSubscription string, message string) error {
	self.Logger.Debugf("Push sent to: %s: %s", receiverSubscription, message)
	return nil
}

func (self *PushService) sendReal(receiverSubscription string, message string) error {
	return nil
}

func (self *PushService) Close(ctx context.Context) error {
	self.Logger.Info("Closing push service")

	return nil
}
