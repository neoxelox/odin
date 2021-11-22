package user

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aodin/date"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
)

type UpdaterUsecase struct {
	class.Usecase
	userRepository repository.UserRepository
}

func NewUpdaterUsecase(configuration internal.Configuration, logger core.Logger, userRepository repository.UserRepository) *UpdaterUsecase {
	return &UpdaterUsecase{
		Usecase:        *class.NewUsecase(configuration, logger),
		userRepository: userRepository,
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

	if updatedUser.Name == "" {
		return nil, ErrInvalidName()
	}

	if updatedUser.Picture == "" {
		picture := fmt.Sprintf(model.USER_DEFAULT_PICTURE, url.QueryEscape(updatedUser.Name))
		updatedUser.Picture = picture
	}

	if user.Name != updatedUser.Name || user.Picture != updatedUser.Picture || !user.Birthday.Equals(updatedUser.Birthday) {
		err := self.userRepository.UpdateProfile(ctx, updatedUser.ID, updatedUser.Name, updatedUser.Picture, updatedUser.Birthday)
		if err != nil {
			return nil, ErrGeneric().Wrap(err)
		}
	}

	return updatedUser, nil
}

func (self *UpdaterUsecase) UpdateEmail(ctx context.Context, user model.User, email string) (*model.User, error) {
	return nil, nil
}

func (self *UpdaterUsecase) UpdatePhone(ctx context.Context, user model.User, email string) (*model.User, error) {
	return nil, nil
}
