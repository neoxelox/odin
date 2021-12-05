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

type UnvoterUsecase struct {
	class.Usecase
	postRepository       repository.PostRepository
	membershipRepository repository.MembershipRepository
}

func NewUnvoterUsecase(configuration internal.Configuration, logger core.Logger, postRepository repository.PostRepository,
	membershipRepository repository.MembershipRepository) *UnvoterUsecase {
	return &UnvoterUsecase{
		Usecase:              *class.NewUsecase(configuration, logger),
		postRepository:       postRepository,
		membershipRepository: membershipRepository,
	}
}

func (self *UnvoterUsecase) Unvote(ctx context.Context, voter model.User, communityID string, postID string) (*model.Post, *model.PostHistory, error) {
	voterMembership, err := self.membershipRepository.GetByUserAndCommunity(ctx, voter.ID, communityID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if voterMembership == nil || voterMembership.DeletedAt != nil {
		return nil, nil, community.ErrNotBelongs()
	}

	post, history, err := self.postRepository.GetByID(ctx, postID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if post == nil || history == nil {
		return nil, nil, ErrInvalid()
	}

	voterIDs := strset.New(post.VoterIDs...)

	if !voterIDs.Has(voterMembership.ID) {
		return post, history, nil
	}

	voterIDs.Remove(voterMembership.ID)

	post.VoterIDs = voterIDs.List()

	if post.Type == model.PostType.ISSUE {
		(*post.Priority)--

		err = self.postRepository.UpdateVotersAndPriority(ctx, post.ID, post.VoterIDs, *post.Priority)
	} else {
		err = self.postRepository.UpdateVoters(ctx, post.ID, post.VoterIDs)
	}

	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	return post, history, nil
}
