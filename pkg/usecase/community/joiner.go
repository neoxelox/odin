package community

import (
	"context"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
)

type JoinerUsecase struct {
	class.Usecase
	communityRepository  repository.CommunityRepository
	membershipRepository repository.MembershipRepository
}

func NewJoinerUsecase(configuration internal.Configuration, logger core.Logger, communityRepository repository.CommunityRepository,
	membershipRepository repository.MembershipRepository) *JoinerUsecase {
	return &JoinerUsecase{
		Usecase:              *class.NewUsecase(configuration, logger),
		communityRepository:  communityRepository,
		membershipRepository: membershipRepository,
	}
}

func (self *JoinerUsecase) Join(ctx context.Context, user model.User, communityID string, door string, role string) (*model.Membership, error) {
	membership := model.NewMembership()
	membership.UserID = user.ID
	membership.CommunityID = communityID
	membership.Door = door
	membership.Role = role

	if len(membership.Door) < model.MEMBERSHIP_DOOR_MIN_LENGTH || len(membership.Door) > model.MEMBERSHIP_DOOR_MAX_LENGTH {
		return nil, ErrInvalidDoor()
	}

	if !model.MembershipRole.Has(membership.Role) {
		return nil, ErrInvalidRole()
	}

	community, err := self.communityRepository.GetByID(ctx, communityID)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	if community == nil {
		return nil, ErrInvalid()
	}

	existingMembership, err := self.membershipRepository.GetByUserAndCommunity(ctx, membership.UserID, membership.CommunityID)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	if existingMembership != nil {
		if existingMembership.DeletedAt == nil {
			return nil, ErrAlreadyJoined()
		}

		existingMembership.DeletedAt = nil

		err = self.membershipRepository.UpdateDeletedAt(ctx, existingMembership.ID, existingMembership.DeletedAt)
		if err != nil {
			return nil, ErrGeneric().Wrap(err)
		}

		return existingMembership, nil
	}

	membership, err = self.membershipRepository.Create(ctx, *membership)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	return membership, nil
}
