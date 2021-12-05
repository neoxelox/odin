package post

import (
	"context"
	"sort"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/utility"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
	"github.com/neoxelox/odin/pkg/usecase/community"
)

type GetterUsecase struct {
	class.Usecase
	postRepository       repository.PostRepository
	membershipRepository repository.MembershipRepository
}

func NewGetterUsecase(configuration internal.Configuration, logger core.Logger, postRepository repository.PostRepository,
	membershipRepository repository.MembershipRepository) *GetterUsecase {
	return &GetterUsecase{
		Usecase:              *class.NewUsecase(configuration, logger),
		postRepository:       postRepository,
		membershipRepository: membershipRepository,
	}
}

func (self *GetterUsecase) Get(ctx context.Context, requester model.User, communityID string, postID string) (*model.Post, *model.PostHistory, error) {
	requesterMembership, err := self.membershipRepository.GetByUserAndCommunity(ctx, requester.ID, communityID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if requesterMembership == nil || requester.DeletedAt != nil {
		return nil, nil, community.ErrNotBelongs()
	}

	post, history, err := self.postRepository.GetByID(ctx, postID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if post == nil || history == nil {
		return nil, nil, ErrInvalid()
	}

	if post.CreatorID != requesterMembership.ID && post.RecipientIDs != nil {
		if !utility.StringIn(requesterMembership.ID, *post.RecipientIDs) {
			return nil, nil, community.ErrNotPermission()
		}
	}

	return post, history, nil
}

func (self *GetterUsecase) GetHistory(ctx context.Context, requester model.User, communityID string, postID string) ([]model.PostHistory, error) {
	requesterMembership, err := self.membershipRepository.GetByUserAndCommunity(ctx, requester.ID, communityID)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	if requesterMembership == nil || requester.DeletedAt != nil {
		return nil, community.ErrNotBelongs()
	}

	post, _, err := self.postRepository.GetByID(ctx, postID)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	if post == nil {
		return nil, ErrInvalid()
	}

	if post.CreatorID != requesterMembership.ID && post.RecipientIDs != nil {
		if !utility.StringIn(requesterMembership.ID, *post.RecipientIDs) {
			return nil, community.ErrNotPermission()
		}
	}

	history, err := self.postRepository.GetHistory(ctx, postID)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	if len(history) < 1 {
		return nil, ErrInvalid()
	}

	sort.Slice(history, func(i, j int) bool {
		return history[i].CreatedAt.Before(history[j].CreatedAt)
	})

	return history, nil
}

func (self *GetterUsecase) GetThread(ctx context.Context, requester model.User, communityID string, threadID string) ([]model.Post, []model.PostHistory, error) {
	requesterMembership, err := self.membershipRepository.GetByUserAndCommunity(ctx, requester.ID, communityID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if requesterMembership == nil || requester.DeletedAt != nil {
		return nil, nil, community.ErrNotBelongs()
	}

	thread, _, err := self.postRepository.GetByID(ctx, threadID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if thread == nil {
		return nil, nil, ErrInvalid()
	}

	if thread.CreatorID != requesterMembership.ID && thread.RecipientIDs != nil {
		if !utility.StringIn(requesterMembership.ID, *thread.RecipientIDs) {
			return nil, nil, community.ErrNotPermission()
		}
	}

	posts, histories, err := self.postRepository.ListByThreadID(ctx, threadID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.Before(posts[j].CreatedAt)
	})

	utility.EqualSort(posts, histories, func(i, j int) bool {
		return *posts[i].LastHistoryID == histories[j].ID
	})

	return posts, histories, nil
}

func (self *GetterUsecase) List(ctx context.Context, requester model.User, communityID string, typee *string) ([]model.Post, []model.PostHistory, error) {
	if typee != nil {
		if !model.PostType.Has(*typee) {
			return nil, nil, ErrInvalidType()
		}
	}

	requesterMembership, err := self.membershipRepository.GetByUserAndCommunity(ctx, requester.ID, communityID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if requesterMembership == nil || requester.DeletedAt != nil {
		return nil, nil, community.ErrNotBelongs()
	}

	posts, histories, err := self.postRepository.ListByCommunityID(ctx, communityID, typee)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.Before(posts[j].CreatedAt)
	})

	utility.EqualSort(posts, histories, func(i, j int) bool {
		return *posts[i].LastHistoryID == histories[j].ID
	})

	return posts, histories, nil
}
