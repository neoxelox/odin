package post

import (
	"context"
	"net/url"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/internal/utility"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
	"github.com/neoxelox/odin/pkg/usecase/community"
)

type CreatorUsecase struct {
	class.Usecase
	database             database.Database
	postRepository       repository.PostRepository
	membershipRepository repository.MembershipRepository
}

func NewCreatorUsecase(configuration internal.Configuration, logger core.Logger, database database.Database,
	postRepository repository.PostRepository, membershipRepository repository.MembershipRepository) *CreatorUsecase {
	return &CreatorUsecase{
		Usecase:              *class.NewUsecase(configuration, logger),
		database:             database,
		postRepository:       postRepository,
		membershipRepository: membershipRepository,
	}
}

func (self *CreatorUsecase) Create(ctx context.Context, creator model.User, communityID string, typee string, threadID *string, priority *int,
	recipientIDs *[]string, message string, categories *[]string, state *string, media *[]string, pollWidgetOptions *[]string) (*model.Post, *model.PostHistory, error) {
	creatorMembership, err := self.membershipRepository.GetByUserAndCommunity(ctx, creator.ID, communityID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if creatorMembership == nil || creatorMembership.DeletedAt != nil {
		return nil, nil, community.ErrNotBelongs()
	}

	if len(message) < model.POST_MESSAGE_MIN_LENGTH || len(message) > model.POST_MESSAGE_MAX_LENGTH {
		return nil, nil, ErrInvalidMessage()
	}

	if media == nil {
		media = &[]string{}
	}

	for i := range *media {
		mediaURL, err := url.ParseRequestURI((*media)[i])
		if err != nil {
			return nil, nil, ErrInvalidMedia().Wrap(err)
		}
		(*media)[i] = mediaURL.String()
	}

	if recipientIDs != nil {
		for _, recipientID := range *recipientIDs {
			if recipientID == creatorMembership.ID {
				return nil, nil, ErrInvalidRecipients()
			}
		}

		recipientMemberships, err := self.membershipRepository.GetByIDsAndCommunity(ctx, *recipientIDs, communityID)
		if err != nil {
			return nil, nil, ErrGeneric().Wrap(err)
		}

		if len(recipientMemberships) != len(*recipientIDs) {
			return nil, nil, ErrInvalidRecipients()
		}
	}

	var post *model.Post
	var history *model.PostHistory

	switch typee {
	case model.PostType.PUBLICATION:
		post, history, err = self.createPublication(ctx, *creatorMembership, threadID, recipientIDs, message, *media, pollWidgetOptions)
		if err != nil {
			return nil, nil, ErrGeneric().As(err)
		}
	case model.PostType.ISSUE:
		post, history, err = self.createIssue(ctx, *creatorMembership, priority, recipientIDs, message, categories, state, *media)
		if err != nil {
			return nil, nil, ErrGeneric().As(err)
		}
	default:
		return nil, nil, ErrInvalidType()
	}

	return post, history, nil
}

func (self *CreatorUsecase) createPublication(ctx context.Context, creator model.Membership, threadID *string, recipientIDs *[]string, message string,
	media []string, pollWidgetOptions *[]string) (*model.Post, *model.PostHistory, error) {
	if threadID != nil {
		thread, _, err := self.postRepository.GetByID(ctx, *threadID)
		if err != nil {
			return nil, nil, ErrGeneric().Wrap(err)
		}

		if thread == nil {
			return nil, nil, ErrInvalidThread()
		}

		if !utility.EqualStringSlice(thread.RecipientIDs, recipientIDs) {
			return nil, nil, ErrInvalidRecipients()
		}

		if thread.CreatorID != creator.ID && thread.RecipientIDs != nil {
			if !utility.StringIn(creator.ID, *thread.RecipientIDs) {
				return nil, nil, community.ErrNotPermission()
			}
		}
	}

	var pollWidget map[string][]string
	if pollWidgetOptions != nil {
		pollWidget = make(map[string][]string, len(*pollWidgetOptions))

		for _, option := range *pollWidgetOptions {
			if len(option) < model.POST_POLL_WIDGET_MIN_OPTION_LENGTH || len(option) > model.POST_POLL_WIDGET_MAX_OPTION_LENGTH {
				return nil, nil, ErrInvalidPoll()
			}
			pollWidget[option] = []string{}
		}

		if len(pollWidget) < model.POST_POLL_WIDGET_MIN_OPTIONS || len(pollWidget) > model.POST_POLL_WIDGET_MAX_OPTIONS {
			return nil, nil, ErrInvalidPoll()
		}
	}

	post := model.NewPost()
	history := model.NewPostHistory()

	post.ThreadID = threadID
	post.CreatorID = creator.ID
	post.LastHistoryID = &history.ID
	post.Type = model.PostType.PUBLICATION
	post.RecipientIDs = recipientIDs

	history.PostID = post.ID
	history.Message = message
	history.Media = media
	history.Widgets = model.PostWidgets{
		Poll: &pollWidget,
	}

	err := self.database.Transaction(ctx, func(ctx context.Context) error {
		var err error

		post, err = self.postRepository.CreatePost(ctx, *post)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		history, err = self.postRepository.CreateHistory(ctx, *history)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	return post, history, nil
}

func (self *CreatorUsecase) createIssue(ctx context.Context, creator model.Membership, priority *int, recipientIDs *[]string, message string,
	categories *[]string, state *string, media []string) (*model.Post, *model.PostHistory, error) {
	if priority == nil {
		return nil, nil, ErrInvalidPriority()
	}

	if state == nil {
		state = &model.PostState.PENDING
	}

	if categories == nil {
		categories = &[]string{}
	}

	if !model.PostState.Has(*state) {
		return nil, nil, ErrInvalidState()
	}

	post := model.NewPost()
	history := model.NewPostHistory()

	post.CreatorID = creator.ID
	post.LastHistoryID = &history.ID
	post.Type = model.PostType.ISSUE
	post.Priority = priority
	post.RecipientIDs = recipientIDs

	history.PostID = post.ID
	history.Message = message
	history.Categories = *categories
	history.State = state
	history.Media = media

	err := self.database.Transaction(ctx, func(ctx context.Context) error {
		var err error

		post, err = self.postRepository.CreatePost(ctx, *post)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		history, err = self.postRepository.CreateHistory(ctx, *history)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	return post, history, nil
}
