package user

import (
	"context"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
	"github.com/nyaruka/phonenumbers"
)

type CreatorUsecase struct {
	class.Usecase
	userRepository repository.UserRepository
}

func NewCreatorUsecase(configuration internal.Configuration, logger core.Logger, userRepository repository.UserRepository) *CreatorUsecase {
	return &CreatorUsecase{
		Usecase:        *class.NewUsecase(configuration, logger),
		userRepository: userRepository,
	}
}

func (self *CreatorUsecase) Create(ctx context.Context, phone string) (*model.User, error) {
	ph, err := phonenumbers.Parse(phone, "ES")
	if err != nil {
		return nil, ErrInvalidPhone().Wrap(err)
	}

	if !phonenumbers.IsValidNumber(ph) {
		return nil, ErrInvalidPhone()
	}

	phone = phonenumbers.Format(ph, phonenumbers.E164)

	user := model.NewUser()
	user.Phone = phone

	user, err = self.userRepository.Create(ctx, *user)
	if err != nil {
		if repository.ErrUserExists().Is(err) {
			return nil, ErrPhoneExists().Wrap(err)
		}

		return nil, ErrGeneric().Wrap(err)
	}

	return user, nil
}
