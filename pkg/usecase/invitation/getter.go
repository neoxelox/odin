package invitation

import (
	"context"
	"sort"
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
)

type GetterUsecase struct {
	class.Usecase
	invitationRepository repository.InvitationRepository
}

func NewGetterUsecase(configuration internal.Configuration, logger core.Logger, invitationRepository repository.InvitationRepository) *GetterUsecase {
	return &GetterUsecase{
		Usecase:              *class.NewUsecase(configuration, logger),
		invitationRepository: invitationRepository,
	}
}

func (self *GetterUsecase) List(ctx context.Context, user model.User) ([]model.Invitation, error) {
	now := time.Now()

	invitations, err := self.invitationRepository.List(ctx, user.Phone)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	filteredInvitations := []model.Invitation{}
	expiredInvitationIDs := []string{}

	for _, invitation := range invitations {
		if now.After(invitation.ExpiresAt) {
			expiredInvitationIDs = append(expiredInvitationIDs, invitation.ID)
		} else {
			filteredInvitations = append(filteredInvitations, invitation)
		}
	}

	if len(expiredInvitationIDs) > 0 {
		err = self.invitationRepository.DeleteByIDs(ctx, expiredInvitationIDs)
		if err != nil {
			return nil, ErrGeneric().Wrap(err)
		}
	}

	sort.Slice(filteredInvitations, func(i, j int) bool {
		return filteredInvitations[i].CreatedAt.Before(filteredInvitations[j].CreatedAt)
	})

	return filteredInvitations, nil
}
