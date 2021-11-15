package view

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/payload"
)

type LoginView struct {
	class.View
}

func NewLoginView(configuration internal.Configuration, logger core.Logger) *LoginView {
	return &LoginView{
		View: *class.NewView(configuration, logger),
	}
}

func (self *LoginView) PostStart() (*payload.PostLoginStartRequest, func(ctx echo.Context) error) {
	request := &payload.PostLoginStartRequest{}
	return request, func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, fmt.Sprintf("%s\n", request.Phone))
	}
}
