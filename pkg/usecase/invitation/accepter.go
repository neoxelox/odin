package invitation

import (
	"context"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
	"github.com/neoxelox/odin/pkg/usecase/community"
)

type AccepterUsecase struct {
	class.Usecase
	database             database.Database
	communityJoiner      community.JoinerUsecase
	invitationRepository repository.InvitationRepository
}

func NewAccepterUsecase(configuration internal.Configuration, logger core.Logger, database database.Database,
	communityJoiner community.JoinerUsecase, invitationRepository repository.InvitationRepository) *AccepterUsecase {
	return &AccepterUsecase{
		Usecase:              *class.NewUsecase(configuration, logger),
		database:             database,
		communityJoiner:      communityJoiner,
		invitationRepository: invitationRepository,
	}
}

func (self *AccepterUsecase) Accept(ctx context.Context, user model.User, invitationID string) (*model.Membership, error) {
	invitation, err := self.invitationRepository.GetByID(ctx, invitationID)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	if invitation == nil {
		return nil, ErrInvalid()
	}

	if invitation.Phone != user.Phone {
		return nil, ErrInvalid()
	}

	var membership *model.Membership
	err = self.database.Transaction(ctx, func(ctx context.Context) error {
		err := self.invitationRepository.DeleteByID(ctx, invitation.ID)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		membership, err = self.communityJoiner.Join(ctx, user, invitation.CommunityID, invitation.Door, invitation.Role)
		if err != nil {
			return ErrGeneric().As(err)
		}

		return nil
	})
	if err != nil {
		return nil, ErrGeneric().As(err)
	}

	return membership, nil
}
