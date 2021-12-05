package post

import (
	"context"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/utility"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
	"github.com/neoxelox/odin/pkg/usecase/community"
	"github.com/scylladb/go-set/strset"
)

type PollVoterUsecase struct {
	class.Usecase
	postRepository       repository.PostRepository
	membershipRepository repository.MembershipRepository
}

func NewPollVoterUsecase(configuration internal.Configuration, logger core.Logger, postRepository repository.PostRepository,
	membershipRepository repository.MembershipRepository) *PollVoterUsecase {
	return &PollVoterUsecase{
		Usecase:              *class.NewUsecase(configuration, logger),
		postRepository:       postRepository,
		membershipRepository: membershipRepository,
	}
}

func (self *PollVoterUsecase) Vote(ctx context.Context, voter model.User, communityID string, postID string, pollOption string) (*model.Post, *model.PostHistory, error) {
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

	if history.Widgets.Poll == nil {
		return nil, nil, ErrInvalidPoll()
	}

	optionVoterIDs, exists := (*history.Widgets.Poll)[pollOption]

	if !exists {
		return nil, nil, ErrInvalidPoll()
	}

	for option := range *history.Widgets.Poll {
		if option != pollOption && utility.StringIn(voterMembership.ID, (*history.Widgets.Poll)[option]) {
			return nil, nil, ErrAlreadyVoted()
		}
	}

	voterIDs := strset.New(optionVoterIDs...)

	if voterIDs.Has(voterMembership.ID) {
		return post, history, nil
	}

	voterIDs.Add(voterMembership.ID)

	(*history.Widgets.Poll)[pollOption] = voterIDs.List()

	err = self.postRepository.UpdateWidgets(ctx, history.ID, history.Widgets)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	return post, history, nil
}
