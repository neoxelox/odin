package post

import (
	"context"
	"net/url"
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/internal/utility"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
	"github.com/neoxelox/odin/pkg/usecase/community"
	"github.com/rs/xid"
)

type UpdaterUsecase struct {
	class.Usecase
	database             database.Database
	postRepository       repository.PostRepository
	membershipRepository repository.MembershipRepository
}

func NewUpdaterUsecase(configuration internal.Configuration, logger core.Logger, database database.Database,
	postRepository repository.PostRepository, membershipRepository repository.MembershipRepository) *UpdaterUsecase {
	return &UpdaterUsecase{
		Usecase:              *class.NewUsecase(configuration, logger),
		database:             database,
		postRepository:       postRepository,
		membershipRepository: membershipRepository,
	}
}

func (self *UpdaterUsecase) Update(ctx context.Context, updator model.User, communityID string, postID string, message *string, categories *[]string,
	state *string, media *[]string, pollWidgetOptions *[]string) (*model.Post, *model.PostHistory, error) {
	updatorMembership, err := self.membershipRepository.GetByUserAndCommunity(ctx, updator.ID, communityID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if updatorMembership == nil || updatorMembership.DeletedAt != nil {
		return nil, nil, community.ErrNotBelongs()
	}

	post, history, err := self.postRepository.GetByID(ctx, postID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if post == nil || history == nil {
		return nil, nil, ErrInvalid()
	}

	if updatorMembership.ID != post.CreatorID &&
		updatorMembership.Role != model.MembershipRole.ADMINISTRATOR &&
		updatorMembership.Role != model.MembershipRole.PRESIDENT &&
		updatorMembership.Role != model.MembershipRole.SECRETARY {
		return nil, nil, community.ErrNotPermission()
	} else if state != nil {
		return nil, nil, community.ErrNotPermission()
	}

	// TODO: Refactorize this sh*t...
	// Just use validators bro, instead of repeating logic
	changed := false

	updatedHistory := history.Copy()
	updatedHistory.ID = xid.New().String()
	updatedHistory.CreatedAt = time.Now()
	updatedHistory.UpdatorID = updatorMembership.ID

	if message != nil {
		updatedHistory.Message = *message
		changed = changed || history.Message != updatedHistory.Message
	}

	if len(updatedHistory.Message) < model.POST_MESSAGE_MIN_LENGTH || len(updatedHistory.Message) > model.POST_MESSAGE_MAX_LENGTH {
		return nil, nil, ErrInvalidMessage()
	}

	if media != nil {
		updatedHistory.Media = *media
		changed = changed || !utility.EqualStringSlice(&history.Media, &updatedHistory.Media)
	}

	for i := range updatedHistory.Media {
		mediaURL, err := url.ParseRequestURI((updatedHistory.Media)[i])
		if err != nil {
			return nil, nil, ErrInvalidMedia().Wrap(err)
		}
		(updatedHistory.Media)[i] = mediaURL.String()
	}

	switch post.Type {
	case model.PostType.PUBLICATION:
		post, history, err = self.updatePublication(ctx, changed, *post, *updatedHistory, pollWidgetOptions)
		if err != nil {
			return nil, nil, ErrGeneric().As(err)
		}
	case model.PostType.ISSUE:
		post, history, err = self.updateIssue(ctx, changed, *post, *updatedHistory, categories, state)
		if err != nil {
			return nil, nil, ErrGeneric().As(err)
		}
	default:
		return nil, nil, ErrInvalidType()
	}

	return post, history, nil
}

func (self *UpdaterUsecase) updatePublication(ctx context.Context, changed bool, post model.Post, history model.PostHistory, pollWidgetOptions *[]string) (*model.Post, *model.PostHistory, error) {
	if pollWidgetOptions != nil {
		pollWidget := make(map[string][]string, len(*pollWidgetOptions))

		for _, option := range *pollWidgetOptions {
			if len(option) < model.POST_POLL_WIDGET_MIN_OPTION_LENGTH || len(option) > model.POST_POLL_WIDGET_MAX_OPTION_LENGTH {
				return nil, nil, ErrInvalidPoll()
			}
			pollWidget[option] = []string{}
		}

		if len(pollWidget) < model.POST_POLL_WIDGET_MIN_OPTIONS || len(pollWidget) > model.POST_POLL_WIDGET_MAX_OPTIONS {
			return nil, nil, ErrInvalidPoll()
		}

		if history.Widgets.Poll != nil {
			lastPollOptions := []string{}
			for option := range *history.Widgets.Poll {
				lastPollOptions = append(lastPollOptions, option)
			}

			changed = changed || !utility.EqualStringSlice(&lastPollOptions, pollWidgetOptions)
		}

		history.Widgets.Poll = &pollWidget
	}

	if !changed {
		return &post, &history, nil
	}

	post.LastHistoryID = &history.ID

	err := self.database.Transaction(ctx, func(ctx context.Context) error {
		var err error

		err = self.postRepository.UpdateHistory(ctx, post.ID, *post.LastHistoryID)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		newhistory, err := self.postRepository.CreateHistory(ctx, history)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		history = *newhistory

		return nil
	})
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	return &post, &history, nil
}

func (self *UpdaterUsecase) updateIssue(ctx context.Context, changed bool, post model.Post, history model.PostHistory, categories *[]string, state *string) (*model.Post, *model.PostHistory, error) {
	if categories != nil {
		changed = changed || !utility.EqualStringSlice(&history.Categories, categories)
		history.Categories = *categories
	}

	if state != nil {
		changed = changed || history.State != state
		history.State = state
	}

	if !model.PostState.Has(*history.State) {
		return nil, nil, ErrInvalidState()
	}

	if !changed {
		return &post, &history, nil
	}

	post.LastHistoryID = &history.ID

	err := self.database.Transaction(ctx, func(ctx context.Context) error {
		var err error

		err = self.postRepository.UpdateHistory(ctx, post.ID, *post.LastHistoryID)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		newhistory, err := self.postRepository.CreateHistory(ctx, history)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		history = *newhistory

		return nil
	})
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	return &post, &history, nil
}
