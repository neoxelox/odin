package core

import (
	"github.com/labstack/echo/v4"

	"github.com/neoxelox/odin/internal"
)

var ErrBinderGeneric = internal.NewError("Binder failed")

type Binder struct {
	configuration internal.Configuration
	logger        Logger
	binder        echo.DefaultBinder
}

func NewBinder(configuration internal.Configuration, logger Logger) *Binder {
	logger.SetLogger(logger.Logger().With().Str("layer", "binder").Logger())

	return &Binder{
		configuration: configuration,
		logger:        logger,
		binder:        echo.DefaultBinder{},
	}
}

func (self *Binder) Bind(i interface{}, c echo.Context) error {
	err := self.binder.Bind(i, c)
	if err != nil {
		return ErrBinderGeneric().Wrap(err)
	}

	return nil
}
