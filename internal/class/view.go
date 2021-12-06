package class

import (
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/core"

	"github.com/labstack/echo/v4"
)

type process interface {
	Process() error
}

type View struct {
	Configuration internal.Configuration
	Logger        core.Logger
}

func NewView(configuration internal.Configuration, logger core.Logger) *View {
	logger.SetLogger(logger.Logger().With().Str("layer", "view").Logger())

	return &View{
		Configuration: configuration,
		Logger:        logger,
	}
}

type ViewHandler func() error

type KeyFunc func(ctx echo.Context) string

type Endpoint struct {
	Request       interface{}
	CacheKey      KeyFunc
	CacheTTL      time.Duration
	RatelimitKey  KeyFunc
	RatelimitRate string
}

func (self *View) Handle(ctx echo.Context, endpoint Endpoint, handler ViewHandler) error {
	if endpoint.Request != nil {
		err := ctx.Bind(endpoint.Request)
		if err != nil {
			return internal.ExcInvalidRequest.Cause(err)
		}

		err = ctx.Validate(endpoint.Request)
		if err != nil {
			return internal.ExcInvalidRequest.Cause(err)
		}

		v, ok := endpoint.Request.(process)
		if ok {
			err = v.Process()
			if err != nil {
				return err // Let Process() raise its own exceptions
			}
		}
	}

	return handler()
}
