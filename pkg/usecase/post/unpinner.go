package post

import (
	"context"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
	"github.com/neoxelox/odin/pkg/usecase/community"
	"github.com/scylladb/go-set/strset"
)

type UnpinnerUsecase struct {
	class.Usecase
	postRepository       repository.PostRepository
	membershipRepository repository.MembershipRepository
	communityRepository  repository.CommunityRepository
}

func NewUnpinnerUsecase(configuration internal.Configuration, logger core.Logger, postRepository repository.PostRepository,
	membershipRepository repository.MembershipRepository, communityRepository repository.CommunityRepository) *UnpinnerUsecase {
	return &UnpinnerUsecase{
		Usecase:              *class.NewUsecase(configuration, logger),
		postRepository:       postRepository,
		membershipRepository: membershipRepository,
		communityRepository:  communityRepository,
	}
}

func (self *UnpinnerUsecase) Unpin(ctx context.Context, pinner model.User, communityID string, postID string) (*model.Community, error) {
	pinnerMembership, err := self.membershipRepository.GetByUserAndCommunity(ctx, pinner.ID, communityID)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	if pinnerMembership == nil || pinnerMembership.DeletedAt != nil {
		return nil, community.ErrNotBelongs()
	}

	if pinnerMembership.Role != model.MembershipRole.ADMINISTRATOR &&
		pinnerMembership.Role != model.MembershipRole.PRESIDENT &&
		pinnerMembership.Role != model.MembershipRole.SECRETARY {
		return nil, community.ErrNotPermission()
	}

	pinnerCommunity, err := self.communityRepository.GetByID(ctx, pinnerMembership.CommunityID)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	if pinnerCommunity == nil {
		return nil, community.ErrInvalid()
	}

	pinnedIDs := strset.New(pinnerCommunity.PinnedIDs...)

	if !pinnedIDs.Has(postID) {
		return pinnerCommunity, nil
	}

	post, history, err := self.postRepository.GetByID(ctx, postID)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	if post == nil || history == nil {
		return nil, ErrInvalid()
	}

	pinnedIDs.Remove(post.ID)

	pinnerCommunity.PinnedIDs = pinnedIDs.List()

	err = self.communityRepository.UpdatePinned(ctx, pinnerCommunity.ID, pinnerCommunity.PinnedIDs)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	return pinnerCommunity, nil
}
