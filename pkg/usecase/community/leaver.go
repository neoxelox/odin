package community

import (
	"context"
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
)

type LeaverUsecase struct {
	class.Usecase
	communityRepository  repository.CommunityRepository
	membershipRepository repository.MembershipRepository
}

func NewLeaverUsecase(configuration internal.Configuration, logger core.Logger, communityRepository repository.CommunityRepository,
	membershipRepository repository.MembershipRepository) *LeaverUsecase {
	return &LeaverUsecase{
		Usecase:              *class.NewUsecase(configuration, logger),
		communityRepository:  communityRepository,
		membershipRepository: membershipRepository,
	}
}

func (self *LeaverUsecase) Leave(ctx context.Context, user model.User, communityID string) error {
	community, err := self.communityRepository.GetByID(ctx, communityID)
	if err != nil {
		return ErrGeneric().Wrap(err)
	}

	if community == nil {
		return ErrInvalid()
	}

	membership, err := self.membershipRepository.GetByUserAndCommunity(ctx, user.ID, community.ID)
	if err != nil {
		return ErrGeneric().Wrap(err)
	}

	if membership == nil || membership.DeletedAt != nil {
		return ErrNotBelongs()
	}

	now := time.Now()
	err = self.membershipRepository.UpdateDeletedAt(ctx, membership.ID, &now)
	if err != nil {
		return ErrGeneric().Wrap(err)
	}

	return nil
}
