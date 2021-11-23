package view

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/payload"
	"github.com/neoxelox/odin/pkg/usecase/otp"
	"github.com/neoxelox/odin/pkg/usecase/user"
)

type UserView struct {
	class.View
	userGetter  user.GetterUsecase
	userUpdater user.UpdaterUsecase
	otpCreator  otp.CreatorUsecase
}

func NewUserView(configuration internal.Configuration, logger core.Logger, userGetter user.GetterUsecase,
	userUpdater user.UpdaterUsecase, otpCreator otp.CreatorUsecase) *UserView {
	return &UserView{
		View:        *class.NewView(configuration, logger),
		userGetter:  userGetter,
		userUpdater: userUpdater,
		otpCreator:  otpCreator,
	}
}

func (self *UserView) GetProfile() (interface{}, func(ctx echo.Context) error) {
	response := &payload.PostUserProfileResponse{}
	return nil, func(ctx echo.Context) error {
		reqUser := RequestUser(ctx)

		response.User = payload.User{
			ID:       reqUser.ID,
			Phone:    reqUser.Phone,
			Name:     reqUser.Name,
			Email:    reqUser.Email,
			Picture:  reqUser.Picture,
			Birthday: reqUser.Birthday,
		}
		return ctx.JSON(http.StatusOK, response)
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

func (self *UserView) PostEmailStart() (*payload.PostUserEmailStartRequest, func(ctx echo.Context) error) {
	request := &payload.PostUserEmailStartRequest{}
	response := &payload.PostUserEmailStartResponse{}
	return request, func(ctx echo.Context) error {
		newOTP, err := self.otpCreator.Create(ctx.Request().Context(), request.Email, model.OTPType.EMAIL)
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

func (self *UserView) PostEmailEnd() (*payload.PostUserEmailEndRequest, func(ctx echo.Context) error) {
	request := &payload.PostUserEmailEndRequest{}
	response := &payload.PostUserEmailEndResponse{}
	return request, func(ctx echo.Context) error {
		reqUser := RequestUser(ctx)

		email, err := self.userUpdater.UpdateEmail(ctx.Request().Context(), *reqUser, request.ID, request.Code)
		switch {
		case err == nil:
			response.Email = email
			return ctx.JSON(http.StatusOK, response)
		case otp.ErrInvalidOTP().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case otp.ErrMaxAttempts().Is(err):
			return ExcOTPMaxAttempts.Cause(err)
		case otp.ErrWrongCode().Is(err):
			return ExcOTPWrongCode.Cause(err)
		case user.ErrInvalidEmail().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	}
}

func (self *UserView) PostPhoneStart() (*payload.PostUserPhoneStartRequest, func(ctx echo.Context) error) {
	request := &payload.PostUserPhoneStartRequest{}
	response := &payload.PostUserPhoneStartResponse{}
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

func (self *UserView) PostPhoneEnd() (*payload.PostUserPhoneEndRequest, func(ctx echo.Context) error) {
	request := &payload.PostUserPhoneEndRequest{}
	response := &payload.PostUserPhoneEndResponse{}
	return request, func(ctx echo.Context) error {
		reqUser := RequestUser(ctx)

		phone, err := self.userUpdater.UpdatePhone(ctx.Request().Context(), *reqUser, request.ID, request.Code)
		switch {
		case err == nil:
			response.Phone = phone
			return ctx.JSON(http.StatusOK, response)
		case otp.ErrInvalidOTP().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case otp.ErrMaxAttempts().Is(err):
			return ExcOTPMaxAttempts.Cause(err)
		case otp.ErrWrongCode().Is(err):
			return ExcOTPWrongCode.Cause(err)
		case user.ErrInvalidPhone().Is(err), user.ErrPhoneExists().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	}
}
