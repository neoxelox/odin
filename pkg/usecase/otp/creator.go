package otp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/internal/utility"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
	"github.com/neoxelox/odin/pkg/service"
)

type CreatorUsecase struct {
	class.Usecase
	database      database.Database
	otpRepository repository.OTPRepository
	smsService    service.SMSService
}

func NewCreatorUsecase(configuration internal.Configuration, logger core.Logger, database database.Database,
	otpRepository repository.OTPRepository, smsService service.SMSService) *CreatorUsecase {
	return &CreatorUsecase{
		Usecase:       *class.NewUsecase(configuration, logger),
		database:      database,
		otpRepository: otpRepository,
		smsService:    smsService,
	}
}

func (self *CreatorUsecase) Create(ctx context.Context, asset string, typee string) (*model.OTP, error) {
	if !model.OTPType.Has(typee) {
		return nil, ErrGeneric()
	}

	existingOTP, err := self.otpRepository.GetByAsset(ctx, asset)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	if existingOTP != nil {
		if time.Now().Before(existingOTP.ExpiresAt) {
			return nil, ErrAlreadySend()
		}

		err = self.otpRepository.Delete(ctx, existingOTP.ID)
		if err != nil {
			return nil, ErrGeneric().Wrap(err)
		}
	}

	otp := model.NewOTP()
	otp.Asset = asset
	otp.Type = typee

	otp.Code = strings.ToUpper(*utility.RandomString(model.OTP_CODE_LENGTH))
	if self.Configuration.Environment == internal.Environment.DEVELOPMENT {
		otp.Code = "123456"
	}

	err = self.database.Transaction(ctx, func(ctx context.Context) error {
		otp, err = self.otpRepository.Create(ctx, *otp)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		err = self.send(ctx, *otp)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		return nil
	})
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	return otp, nil
}

func (self *CreatorUsecase) send(ctx context.Context, otp model.OTP) error {
	message := fmt.Sprintf(OTP_MESSAGE, otp.Code)

	switch otp.Type {
	case model.OTPType.SMS:
		err := self.smsService.Send(otp.Asset, message)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}
		return nil
	case model.OTPType.EMAIL:
		return ErrGeneric()
	default:
		return ErrGeneric()
	}
}
