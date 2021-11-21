package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/usecase/auth"
)

const (
	AUTH_HEADER = "Authorization"
)

type AuthMiddleware struct {
	class.Middleware
	authVerifier auth.VerifierUsecase
}

func NewAuthMiddleware(configuration internal.Configuration, logger core.Logger, authVerifier auth.VerifierUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		Middleware:   *class.NewMiddleware(configuration, logger),
		authVerifier: authVerifier,
	}
}

func (self *AuthMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if echoMiddleware.DefaultSkipper(ctx) {
			return next(ctx)
		}

		user, err := self.authVerifier.Verify(ctx.Request().Context(), ctx.Request().Header.Get(AUTH_HEADER), model.SessionMetadata{
			IP:         ctx.RealIP(),
			Device:     ctx.Request().Host,
			ApiVersion: strings.Split(ctx.Path(), "/")[1],
		})

		switch {
		case err == nil:
			ctx.Set(string(model.CONTEXT_USER_KEY), user)
			return next(ctx)
		case auth.ErrExpiredAccessToken().Is(err), auth.ErrInvalidAccessToken().Is(err),
			auth.ErrTamperedAccessToken().Is(err), auth.ErrUserBanned().Is(err):
			return internal.ExcUnauthorized.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	}
}
