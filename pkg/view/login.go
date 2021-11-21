package view

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/payload"
	"github.com/neoxelox/odin/pkg/usecase/login"
	"github.com/neoxelox/odin/pkg/usecase/otp"
)

type LoginView struct {
	class.View
	otpCreator     otp.CreatorUsecase
	loginProcessor login.ProcessorUsecase
}

func NewLoginView(configuration internal.Configuration, logger core.Logger, otpCreator otp.CreatorUsecase,
	loginProcessor login.ProcessorUsecase) *LoginView {
	return &LoginView{
		View:           *class.NewView(configuration, logger),
		otpCreator:     otpCreator,
		loginProcessor: loginProcessor,
	}
}

func (self *LoginView) PostStart() (*payload.PostLoginStartRequest, func(ctx echo.Context) error) {
	request := &payload.PostLoginStartRequest{}
	response := &payload.PostLoginStartResponse{}
	return request, func(ctx echo.Context) error {
		newOTP, err := self.otpCreator.Create(ctx.Request().Context(), request.Phone, model.OTPType.SMS)
		switch {
		case err == nil:
			response.ID = newOTP.ID
			return ctx.JSON(http.StatusOK, response)
		case otp.ErrAlreadySend().Is(err):
			return ExcOTPAlreadySend.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	}
}

func (self *LoginView) PostEnd() (*payload.PostLoginEndRequest, func(ctx echo.Context) error) {
	request := &payload.PostLoginEndRequest{}
	response := &payload.PostLoginEndResponse{}
	return request, func(ctx echo.Context) error {
		accessToken, user, err := self.loginProcessor.Process(ctx.Request().Context(), request.ID, request.Code, model.SessionMetadata{
			IP:         ctx.RealIP(),
			Device:     ctx.Request().Host,
			ApiVersion: "v1",
		})
		switch {
		case err == nil:
			response.AccessToken = accessToken
			if user != nil {
				response.User = &payload.User{
					ID:       user.ID,
					Phone:    user.Phone,
					Name:     user.Name,
					Email:    user.Email,
					Picture:  user.Picture,
					Birthday: user.Birthday,
				}
			}
			return ctx.JSON(http.StatusOK, response)
		case otp.ErrInvalidOTP().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case otp.ErrMaxAttempts().Is(err):
			return ExcOTPMaxAttempts.Cause(err)
		case otp.ErrWrongCode().Is(err):
			return ExcOTPWrongCode.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	}
}
