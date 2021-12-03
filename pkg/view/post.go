package view

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/payload"
	"github.com/neoxelox/odin/pkg/usecase/post"
)

type PostView struct {
	class.View
	postCreator post.CreatorUsecase
}

func NewPostView(configuration internal.Configuration, logger core.Logger, postCreator post.CreatorUsecase) *PostView {
	return &PostView{
		View:        *class.NewView(configuration, logger),
		postCreator: postCreator,
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
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}
