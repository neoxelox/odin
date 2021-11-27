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
	userDeleter user.DeleterUsecase
	otpCreator  otp.CreatorUsecase
}

func NewUserView(configuration internal.Configuration, logger core.Logger, userGetter user.GetterUsecase,
	userUpdater user.UpdaterUsecase, userDeleter user.DeleterUsecase, otpCreator otp.CreatorUsecase) *UserView {
	return &UserView{
		View:        *class.NewView(configuration, logger),
		userGetter:  userGetter,
		userUpdater: userUpdater,
		userDeleter: userDeleter,
		otpCreator:  otpCreator,
	}
}

func (self *UserView) GetProfile(ctx echo.Context) error {
	requestUser := RequestUser(ctx)
	response := &payload.PostUserProfileResponse{}
	return self.Handle(ctx, class.Endpoint{}, func() error {
		response.User = payload.User{
			ID:       requestUser.ID,
			Phone:    requestUser.Phone,
			Name:     requestUser.Name,
			Email:    requestUser.Email,
			Picture:  requestUser.Picture,
			Birthday: requestUser.Birthday,
		}

		return ctx.JSON(http.StatusOK, response)
	})
}

func (self *UserView) PostProfile(ctx echo.Context) error {
	request := &payload.PostUserProfileRequest{}
	requestUser := RequestUser(ctx)
	response := &payload.PostUserProfileResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		updatedUser, err := self.userUpdater.UpdateProfile(ctx.Request().Context(), *requestUser, request.Name, request.LastName,
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
	})
}

func (self *UserView) PostEmailStart(ctx echo.Context) error {
	request := &payload.PostUserEmailStartRequest{}
	response := &payload.PostUserEmailStartResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
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
	})
}

func (self *UserView) PostEmailEnd(ctx echo.Context) error {
	request := &payload.PostUserEmailEndRequest{}
	requestUser := RequestUser(ctx)
	response := &payload.PostUserEmailEndResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		email, err := self.userUpdater.UpdateEmail(ctx.Request().Context(), *requestUser, request.ID, request.Code)
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
	})
}

func (self *UserView) PostPhoneStart(ctx echo.Context) error {
	request := &payload.PostUserPhoneStartRequest{}
	response := &payload.PostUserPhoneStartResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
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
	})
}

func (self *UserView) PostPhoneEnd(ctx echo.Context) error {
	request := &payload.PostUserPhoneEndRequest{}
	requestUser := RequestUser(ctx)
	response := &payload.PostUserPhoneEndResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		phone, err := self.userUpdater.UpdatePhone(ctx.Request().Context(), *requestUser, request.ID, request.Code)
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
	})
}

func (self *UserView) DeleteUser(ctx echo.Context) error {
	requestUser := RequestUser(ctx)
	response := &payload.DeleteUserResponse{}
	return self.Handle(ctx, class.Endpoint{}, func() error {
		err := self.userDeleter.Delete(ctx.Request().Context(), *requestUser)
		switch {
		case err == nil:
			return ctx.JSON(http.StatusOK, response)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}
