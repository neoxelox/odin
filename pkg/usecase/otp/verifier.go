package otp

import (
	"context"
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
)

type VerifierUsecase struct {
	class.Usecase
	otpRepository repository.OTPRepository
}

func NewVerifierUsecase(configuration internal.Configuration, logger core.Logger, otpRepository repository.OTPRepository) *VerifierUsecase {
	return &VerifierUsecase{
		Usecase:       *class.NewUsecase(configuration, logger),
		otpRepository: otpRepository,
	}
}

func (self *VerifierUsecase) Verify(ctx context.Context, otpID string, code string, typee string) (*model.OTP, error) {
	if !model.OTPType.Has(typee) {
		return nil, ErrGeneric()
	}

	otp, err := self.otpRepository.GetByID(ctx, otpID)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	if otp == nil {
		return nil, ErrInvalid()
	}

	if otp.Type != typee {
		return nil, ErrInvalid()
	}

	if time.Now().After(otp.ExpiresAt) {
		return nil, ErrInvalid()
	}

	otp.Attempts++

	if otp.Attempts > model.OTP_MAX_ATTEMPTS {
		return nil, ErrMaxAttempts()
	}

	err = self.otpRepository.UpdateAttempts(ctx, otp.ID, otp.Attempts)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	if otp.Code != code {
		return nil, ErrWrongCode()
	}

	return otp, nil
}
