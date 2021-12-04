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
	renderer      core.Renderer
	otpRepository repository.OTPRepository
	smsService    service.SMSService
	emailService  service.EmailService
}

func NewCreatorUsecase(configuration internal.Configuration, logger core.Logger, database database.Database,
	renderer core.Renderer, otpRepository repository.OTPRepository, smsService service.SMSService,
	emailService service.EmailService) *CreatorUsecase {
	return &CreatorUsecase{
		Usecase:       *class.NewUsecase(configuration, logger),
		database:      database,
		renderer:      renderer,
		otpRepository: otpRepository,
		smsService:    smsService,
		emailService:  emailService,
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
			return nil, ErrAlreadySent()
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
	} else { // TODO: Temporary DELETE!!
		otp.Code = "DM5FJH"
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
	switch otp.Type {
	case model.OTPType.SMS:
		err := self.smsService.Send(otp.Asset, fmt.Sprintf(OTP_SMS_MESSAGE, otp.Code))
		if err != nil {
			return ErrGeneric().Wrap(err)
		}
		return nil
	case model.OTPType.EMAIL:
		body, err := self.renderer.RenderString(OTP_EMAIL_TEMPLATE, otp.Code)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}
		err = self.emailService.Send(otp.Asset, OTP_EMAIL_SUBJECT, body)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}
		return nil
	default:
		return ErrGeneric()
	}
}
