package user

import (
	"context"
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
)

type DeleterUsecase struct {
	class.Usecase
	userRepository repository.UserRepository
}

func NewDeleterUsecase(configuration internal.Configuration, logger core.Logger, userRepository repository.UserRepository) *DeleterUsecase {
	return &DeleterUsecase{
		Usecase:        *class.NewUsecase(configuration, logger),
		userRepository: userRepository,
	}
}

func (self *DeleterUsecase) Delete(ctx context.Context, user model.User) error {
	err := self.userRepository.UpdateDeletedAt(ctx, user.ID, time.Now())
	if err != nil {
		return ErrGeneric().Wrap(err)
	}

	return nil
}
