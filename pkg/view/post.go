package view

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/payload"
	"github.com/neoxelox/odin/pkg/usecase/community"
	"github.com/neoxelox/odin/pkg/usecase/post"
)

type PostView struct {
	class.View
	postCreator   post.CreatorUsecase
	postUpdater   post.UpdaterUsecase
	postVoter     post.VoterUsecase
	postUnvoter   post.UnvoterUsecase
	postPollVoter post.PollVoterUsecase
	postPinner    post.PinnerUsecase
	postUnpinner  post.UnpinnerUsecase
}

func NewPostView(configuration internal.Configuration, logger core.Logger, postCreator post.CreatorUsecase,
	postUpdater post.UpdaterUsecase, postVoter post.VoterUsecase, postUnvoter post.UnvoterUsecase,
	postPollVoter post.PollVoterUsecase, postPinner post.PinnerUsecase, postUnpinner post.UnpinnerUsecase) *PostView {
	return &PostView{
		View:          *class.NewView(configuration, logger),
		postCreator:   postCreator,
		postUpdater:   postUpdater,
		postVoter:     postVoter,
		postUnvoter:   postUnvoter,
		postPollVoter: postPollVoter,
		postPinner:    postPinner,
		postUnpinner:  postUnpinner,
	}
}

func (self *PostView) PostPost(ctx echo.Context) error {
	request := &payload.PostPostRequest{
		Widgets: &struct { // WTF Golang...?
			PollOptions *[]string "json:\"poll_options\" validate:\"omitempty,required\""
		}{},
	}
	requestUser := RequestUser(ctx)
	response := &payload.PostPostResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		newPost, newHistory, err := self.postCreator.Create(ctx.Request().Context(), *requestUser, request.CommunityID, request.Type, request.ThreadID, request.Priority,
			request.RecipientIDs, request.Message, request.Categories, request.State, request.Media, request.Widgets.PollOptions)
		switch {
		case err == nil:
			response.Post = payload.Post{
				ID:           newPost.ID,
				ThreadID:     newPost.ThreadID,
				CreatorID:    newPost.CreatorID,
				Type:         newPost.Type,
				Priority:     newPost.Priority,
				RecipientIDs: newPost.RecipientIDs,
				VoterIDs:     newPost.VoterIDs,
				CreatedAt:    newPost.CreatedAt,
				PostHistory: payload.PostHistory{
					Message:    newHistory.Message,
					Categories: newHistory.Categories,
					State:      newHistory.State,
					Media:      newHistory.Media,
					Widgets: payload.PostWidgets{
						Poll: newHistory.Widgets.Poll,
					},
				},
			}
			return ctx.JSON(http.StatusOK, response)
		case post.ErrInvalid().Is(err), post.ErrInvalidType().Is(err), post.ErrInvalidThread().Is(err),
			post.ErrInvalidPriority().Is(err), post.ErrInvalidRecipients().Is(err), post.ErrInvalidMessage().Is(err),
			post.ErrInvalidState().Is(err), post.ErrInvalidMedia().Is(err), post.ErrInvalidPoll().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case community.ErrNotBelongs().Is(err):
			return ExcUserNotBelongs.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}

func (self *PostView) PutPost(ctx echo.Context) error {
	request := &payload.PutPostRequest{
		Widgets: &struct { // WTF Golang...?
			PollOptions *[]string "json:\"poll_options\" validate:\"omitempty,required\""
		}{},
	}
	requestUser := RequestUser(ctx)
	response := &payload.PutPostResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		updatedPost, updatedHistory, err := self.postUpdater.Update(ctx.Request().Context(), *requestUser, request.CommunityID, request.PostID, request.Message, request.Categories,
			request.State, request.Media, request.Widgets.PollOptions)
		switch {
		case err == nil:
			response.Post = payload.Post{
				ID:           updatedPost.ID,
				ThreadID:     updatedPost.ThreadID,
				CreatorID:    updatedPost.CreatorID,
				Type:         updatedPost.Type,
				Priority:     updatedPost.Priority,
				RecipientIDs: updatedPost.RecipientIDs,
				VoterIDs:     updatedPost.VoterIDs,
				CreatedAt:    updatedPost.CreatedAt,
				PostHistory: payload.PostHistory{
					Message:    updatedHistory.Message,
					Categories: updatedHistory.Categories,
					State:      updatedHistory.State,
					Media:      updatedHistory.Media,
					Widgets: payload.PostWidgets{
						Poll: updatedHistory.Widgets.Poll,
					},
				},
			}
			return ctx.JSON(http.StatusOK, response)
		case post.ErrInvalid().Is(err), post.ErrInvalidType().Is(err), post.ErrInvalidMessage().Is(err),
			post.ErrInvalidState().Is(err), post.ErrInvalidMedia().Is(err), post.ErrInvalidPoll().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case community.ErrNotBelongs().Is(err):
			return ExcUserNotBelongs.Cause(err)
		case community.ErrNotPermission().Is(err):
			return ExcUserNotPermission.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}

func (self *PostView) PostVotePost(ctx echo.Context) error {
	request := &payload.PostVotePostRequest{}
	requestUser := RequestUser(ctx)
	response := &payload.PostVotePostResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		resPost, resHistory, err := self.postVoter.Vote(ctx.Request().Context(), *requestUser, request.CommunityID, request.PostID)
		switch {
		case err == nil:
			response.Post = payload.Post{
				ID:           resPost.ID,
				ThreadID:     resPost.ThreadID,
				CreatorID:    resPost.CreatorID,
				Type:         resPost.Type,
				Priority:     resPost.Priority,
				RecipientIDs: resPost.RecipientIDs,
				VoterIDs:     resPost.VoterIDs,
				CreatedAt:    resPost.CreatedAt,
				PostHistory: payload.PostHistory{
					Message:    resHistory.Message,
					Categories: resHistory.Categories,
					State:      resHistory.State,
					Media:      resHistory.Media,
					Widgets: payload.PostWidgets{
						Poll: resHistory.Widgets.Poll,
					},
				},
			}
			return ctx.JSON(http.StatusOK, response)
		case post.ErrInvalid().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case community.ErrNotBelongs().Is(err):
			return ExcUserNotBelongs.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}

func (self *PostView) PostUnvotePost(ctx echo.Context) error {
	request := &payload.PostUnvotePostRequest{}
	requestUser := RequestUser(ctx)
	response := &payload.PostUnvotePostResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		resPost, resHistory, err := self.postUnvoter.Unvote(ctx.Request().Context(), *requestUser, request.CommunityID, request.PostID)
		switch {
		case err == nil:
			response.Post = payload.Post{
				ID:           resPost.ID,
				ThreadID:     resPost.ThreadID,
				CreatorID:    resPost.CreatorID,
				Type:         resPost.Type,
				Priority:     resPost.Priority,
				RecipientIDs: resPost.RecipientIDs,
				VoterIDs:     resPost.VoterIDs,
				CreatedAt:    resPost.CreatedAt,
				PostHistory: payload.PostHistory{
					Message:    resHistory.Message,
					Categories: resHistory.Categories,
					State:      resHistory.State,
					Media:      resHistory.Media,
					Widgets: payload.PostWidgets{
						Poll: resHistory.Widgets.Poll,
					},
				},
			}
			return ctx.JSON(http.StatusOK, response)
		case post.ErrInvalid().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case community.ErrNotBelongs().Is(err):
			return ExcUserNotBelongs.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}

func (self *PostView) PostVotePostPoll(ctx echo.Context) error {
	request := &payload.PostVotePostPollRequest{}
	requestUser := RequestUser(ctx)
	response := &payload.PostVotePostPollResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		resPost, resHistory, err := self.postPollVoter.Vote(ctx.Request().Context(), *requestUser, request.CommunityID, request.PostID, request.Option)
		switch {
		case err == nil:
			response.Post = payload.Post{
				ID:           resPost.ID,
				ThreadID:     resPost.ThreadID,
				CreatorID:    resPost.CreatorID,
				Type:         resPost.Type,
				Priority:     resPost.Priority,
				RecipientIDs: resPost.RecipientIDs,
				VoterIDs:     resPost.VoterIDs,
				CreatedAt:    resPost.CreatedAt,
				PostHistory: payload.PostHistory{
					Message:    resHistory.Message,
					Categories: resHistory.Categories,
					State:      resHistory.State,
					Media:      resHistory.Media,
					Widgets: payload.PostWidgets{
						Poll: resHistory.Widgets.Poll,
					},
				},
			}
			return ctx.JSON(http.StatusOK, response)
		case post.ErrInvalid().Is(err), post.ErrInvalidPoll().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case community.ErrNotBelongs().Is(err):
			return ExcUserNotBelongs.Cause(err)
		case post.ErrAlreadyVoted().Is(err):
			return ExcUserAlreadyVoted.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}

func (self *PostView) PostPinPost(ctx echo.Context) error {
	request := &payload.PostPinPostRequest{}
	requestUser := RequestUser(ctx)
	response := &payload.PostPinPostResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		resCommunity, err := self.postPinner.Pin(ctx.Request().Context(), *requestUser, request.CommunityID, request.PostID)
		switch {
		case err == nil:
			response.Community = payload.Community{
				ID:         resCommunity.ID,
				Address:    resCommunity.Address,
				Name:       resCommunity.Name,
				Categories: resCommunity.Categories,
				PinnedIDs:  resCommunity.PinnedIDs,
				CreatedAt:  resCommunity.CreatedAt,
			}
			return ctx.JSON(http.StatusOK, response)
		case community.ErrInvalid().Is(err), post.ErrInvalid().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case community.ErrNotBelongs().Is(err):
			return ExcUserNotBelongs.Cause(err)
		case community.ErrNotPermission().Is(err):
			return ExcUserNotPermission.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}

func (self *PostView) PostUnpinPost(ctx echo.Context) error {
	request := &payload.PostUnpinPostRequest{}
	requestUser := RequestUser(ctx)
	response := &payload.PostUnpinPostResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		resCommunity, err := self.postUnpinner.Unpin(ctx.Request().Context(), *requestUser, request.CommunityID, request.PostID)
		switch {
		case err == nil:
			response.Community = payload.Community{
				ID:         resCommunity.ID,
				Address:    resCommunity.Address,
				Name:       resCommunity.Name,
				Categories: resCommunity.Categories,
				PinnedIDs:  resCommunity.PinnedIDs,
				CreatedAt:  resCommunity.CreatedAt,
			}
			return ctx.JSON(http.StatusOK, response)
		case community.ErrInvalid().Is(err), post.ErrInvalid().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case community.ErrNotBelongs().Is(err):
			return ExcUserNotBelongs.Cause(err)
		case community.ErrNotPermission().Is(err):
			return ExcUserNotPermission.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}
