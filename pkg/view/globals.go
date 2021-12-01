package view

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/pkg/model"
)

var (
	ExcOTPAlreadySent     = internal.NewException(http.StatusForbidden, "ERR_OTP_ALREADY_SENT")
	ExcOTPMaxAttempts     = internal.NewException(http.StatusForbidden, "ERR_OTP_MAX_ATTEMPTS")
	ExcOTPWrongCode       = internal.NewException(http.StatusForbidden, "ERR_OTP_WRONG_CODE")
	ExcUserAlreadyJoined  = internal.NewException(http.StatusForbidden, "ERR_USER_ALREADY_JOINED")
	ExcUserNotBelongs     = internal.NewException(http.StatusForbidden, "ERR_USER_NOT_BELONGS")
	ExcUserNotPermission  = internal.NewException(http.StatusForbidden, "ERR_USER_NOT_PERMISSION")
	ExcUserAlreadyInvited = internal.NewException(http.StatusForbidden, "ERR_USER_ALREADY_INVITED")
)

func RequestSession(ctx echo.Context) *model.Session {
	return ctx.Get(string(model.CONTEXT_SESSION_KEY)).(*model.Session)
}

func RequestUser(ctx echo.Context) *model.User {
	return ctx.Get(string(model.CONTEXT_USER_KEY)).(*model.User)
}
