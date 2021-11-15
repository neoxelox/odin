package class

import (
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/core"
)

type Middleware struct {
	Configuration internal.Configuration
	Logger        core.Logger
}

func NewMiddleware(configuration internal.Configuration, logger core.Logger) *Middleware {
	logger.SetLogger(logger.Logger().With().Str("layer", "middleware").Logger())

	return &Middleware{
		Configuration: configuration,
		Logger:        logger,
	}
}
