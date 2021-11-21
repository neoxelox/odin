package view

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/payload"
)

type UserView struct {
	class.View
}

func NewUserView(configuration internal.Configuration, logger core.Logger) *UserView {
	return &UserView{
		View: *class.NewView(configuration, logger),
	}
}

func (self *UserView) Get() (*payload.GetUserRequest, func(ctx echo.Context) error) {
	request := &payload.GetUserRequest{}
	response := &payload.GetUserResponse{}
	return request, func(ctx echo.Context) error {

		// TODO CHECK IF REQUESTED USER BELONGS TO THE SAME COMMUNITY:

		user := RequestUser(ctx)

		response.User = payload.User{
			ID:       user.ID,
			Phone:    user.Phone,
			Name:     user.Name,
			Email:    user.Email,
			Picture:  user.Picture,
			Birthday: user.Birthday,
		}

		return ctx.JSON(http.StatusOK, response)

	}
}
