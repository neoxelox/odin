package invitation

import (
	"context"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
)

type RejecterUsecase struct {
	class.Usecase
	invitationRepository repository.InvitationRepository
}

func NewRejecterUsecase(configuration internal.Configuration, logger core.Logger, invitationRepository repository.InvitationRepository) *RejecterUsecase {
	return &RejecterUsecase{
		Usecase:              *class.NewUsecase(configuration, logger),
		invitationRepository: invitationRepository,
	}
}

func (self *RejecterUsecase) Reject(ctx context.Context, user model.User, invitationID string) error {
	invitation, err := self.invitationRepository.GetByID(ctx, invitationID)
	if err != nil {
		return ErrGeneric().Wrap(err)
	}

	if invitation == nil {
		return ErrInvalid()
	}

	if invitation.Phone != user.Phone {
		return ErrInvalid()
	}

	err = self.invitationRepository.DeleteByID(ctx, invitation.ID)
	if err != nil {
		return ErrGeneric().Wrap(err)
	}

	return nil
}
