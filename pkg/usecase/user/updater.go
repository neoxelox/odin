package user

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aodin/date"
	"github.com/badoux/checkmail"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
	"github.com/neoxelox/odin/pkg/usecase/otp"
	"github.com/nyaruka/phonenumbers"
)

type UpdaterUsecase struct {
	class.Usecase
	database       database.Database
	userRepository repository.UserRepository
	otpRepository  repository.OTPRepository
	otpVerifier    otp.VerifierUsecase
}

func NewUpdaterUsecase(configuration internal.Configuration, logger core.Logger, database database.Database,
	userRepository repository.UserRepository, otpRepository repository.OTPRepository, otpVerifier otp.VerifierUsecase) *UpdaterUsecase {
	return &UpdaterUsecase{
		Usecase:        *class.NewUsecase(configuration, logger),
		database:       database,
		userRepository: userRepository,
		otpRepository:  otpRepository,
		otpVerifier:    otpVerifier,
	}
}

func (self *UpdaterUsecase) UpdateProfile(ctx context.Context, user model.User, name *string, lastName *string, picture *string,
	birthday *date.Date) (*model.User, error) {
	updatedUser := user.Copy()

	if name != nil && lastName != nil {
		updatedUser.Name = *name + " " + *lastName
	}

	if picture != nil {
		updatedUser.Picture = *picture
	}

	if birthday != nil {
		updatedUser.Birthday = *birthday
	}

	if len(updatedUser.Name) < model.USER_NAME_MIN_LENGTH || len(updatedUser.Name) > model.USER_NAME_MAX_LENGTH {
		return nil, ErrInvalidName()
	}

	if updatedUser.Picture == "" {
		updatedUser.Picture = fmt.Sprintf(model.USER_DEFAULT_PICTURE, url.QueryEscape(updatedUser.Name))
	}

	pictureURL, err := url.ParseRequestURI(updatedUser.Picture)
	if err != nil {
		return nil, ErrInvalidPicture().Wrap(err)
	}
	updatedUser.Picture = pictureURL.String()

	if user.Name != updatedUser.Name || user.Picture != updatedUser.Picture || !user.Birthday.Equals(updatedUser.Birthday) {
		err := self.userRepository.UpdateProfile(ctx, updatedUser.ID, updatedUser.Name, updatedUser.Picture, updatedUser.Birthday)
		if err != nil {
			return nil, ErrGeneric().Wrap(err)
		}
	}

	return updatedUser, nil
}

func (self *UpdaterUsecase) UpdateEmail(ctx context.Context, user model.User, otpID string, code string) (string, error) {
	otpReq, err := self.otpVerifier.Verify(ctx, otpID, code, model.OTPType.EMAIL)
	if err != nil {
		if !otp.ErrGeneric().Is(err) {
			return "", ErrGeneric().As(err)
		}

		return "", ErrGeneric().Wrap(err)
	}

	if otpReq.Asset == user.Email {
		err := self.otpRepository.Delete(ctx, otpReq.ID)
		if err != nil {
			return "", ErrGeneric().Wrap(err)
		}

		return otpReq.Asset, nil
	}

	user.Email = otpReq.Asset

	err = checkmail.ValidateFormat(user.Email)
	if err != nil {
		return "", ErrInvalidEmail().Wrap(err)
	}

	err = self.database.Transaction(ctx, func(ctx context.Context) error {
		err := self.otpRepository.Delete(ctx, otpReq.ID)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		err = self.userRepository.UpdateEmail(ctx, user.ID, user.Email)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		return nil
	})
	if err != nil {
		return "", ErrGeneric().Wrap(err)
	}

	return user.Email, nil
}

func (self *UpdaterUsecase) UpdatePhone(ctx context.Context, user model.User, otpID string, code string) (string, error) {
	otpReq, err := self.otpVerifier.Verify(ctx, otpID, code, model.OTPType.SMS)
	if err != nil {
		if !otp.ErrGeneric().Is(err) {
			return "", ErrGeneric().As(err)
		}

		return "", ErrGeneric().Wrap(err)
	}

	if otpReq.Asset == user.Phone {
		err := self.otpRepository.Delete(ctx, otpReq.ID)
		if err != nil {
			return "", ErrGeneric().Wrap(err)
		}

		return otpReq.Asset, nil
	}

	user.Phone = otpReq.Asset

	ph, err := phonenumbers.Parse(user.Phone, "ES")
	if err != nil {
		return "", ErrInvalidPhone().Wrap(err)
	}

	if !phonenumbers.IsValidNumber(ph) {
		return "", ErrInvalidPhone()
	}

	user.Phone = phonenumbers.Format(ph, phonenumbers.E164)

	err = self.database.Transaction(ctx, func(ctx context.Context) error {
		err := self.otpRepository.Delete(ctx, otpReq.ID)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		err = self.userRepository.UpdatePhone(ctx, user.ID, user.Phone)
		if err != nil {
			if repository.ErrUserExists().Is(err) {
				return ErrPhoneExists().Wrap(err)
			}

			return ErrGeneric().Wrap(err)
		}

		return nil
	})
	if err != nil {
		if !ErrGeneric().Is(err) {
			return "", ErrGeneric().As(err)
		}

		return "", ErrGeneric().Wrap(err)
	}

	return user.Phone, nil
}
