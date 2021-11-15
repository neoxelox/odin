package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
)

type LoggerMiddleware struct {
	class.Middleware
}

func NewLoggerMiddleware(configuration internal.Configuration, logger core.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		Middleware: *class.NewMiddleware(configuration, logger),
	}
}

func (self *LoggerMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if echoMiddleware.DefaultSkipper(ctx) {
			return next(ctx)
		}

		req := ctx.Request()
		res := ctx.Response()

		start := time.Now()

		if err := next(ctx); err != nil {
			ctx.Error(err)
		}

		stop := time.Now()

		self.Logger.Logger().Info().
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Int("status", res.Status).
			Str("ip_address", ctx.RealIP()).
			Dur("latency", stop.Sub(start)).
			Msg("")

		return nil
	}
}
