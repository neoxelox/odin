package user

import (
	"context"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
)

type GetterUsecase struct {
	class.Usecase
	userRepository repository.UserRepository
}

func NewGetterUsecase(configuration internal.Configuration, logger core.Logger, userRepository repository.UserRepository) *GetterUsecase {
	return &GetterUsecase{
		Usecase:        *class.NewUsecase(configuration, logger),
		userRepository: userRepository,
	}
}

func (self *GetterUsecase) Get(ctx context.Context, id string) (*model.User, error) {
	user, err := self.userRepository.GetByID(ctx, id)
	if err != nil || user == nil {
		return nil, ErrGeneric().Wrap(err)
	}

	return user, nil
}
