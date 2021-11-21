package class

import (
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

func (self *View) Handle(request interface{}, handler func(ctx echo.Context) error) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		if request != nil {
			err := ctx.Bind(request)
			if err != nil {
				return internal.ExcInvalidRequest.Cause(err)
			}

			err = ctx.Validate(request)
			if err != nil {
				return internal.ExcInvalidRequest.Cause(err)
			}

			v, ok := request.(process)
			if ok {
				err = v.Process()
				if err != nil {
					return err // Let Process() raise its own exceptions
				}
			}
		}

		return handler(ctx)
	}
}
