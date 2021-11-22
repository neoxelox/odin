package view

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/payload"
	"github.com/neoxelox/odin/pkg/usecase/user"
)

type UserView struct {
	class.View
	userUpdater user.UpdaterUsecase
}

func NewUserView(configuration internal.Configuration, logger core.Logger, userUpdater user.UpdaterUsecase) *UserView {
	return &UserView{
		View:        *class.NewView(configuration, logger),
		userUpdater: userUpdater,
	}
}

func (self *UserView) PostProfile() (*payload.PostUserProfileRequest, func(ctx echo.Context) error) {
	request := &payload.PostUserProfileRequest{}
	response := &payload.PostUserProfileResponse{}
	return request, func(ctx echo.Context) error {
		reqUser := RequestUser(ctx)

		updatedUser, err := self.userUpdater.UpdateProfile(ctx.Request().Context(), *reqUser, request.Name, request.LastName,
			request.Picture, request.Birthday)
		switch {
		case err == nil:
			response.User = payload.User{
				ID:       updatedUser.ID,
				Phone:    updatedUser.Phone,
				Name:     updatedUser.Name,
				Email:    updatedUser.Email,
				Picture:  updatedUser.Picture,
				Birthday: updatedUser.Birthday,
			}
			return ctx.JSON(http.StatusOK, response)
		case user.ErrInvalidName().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	}
}
