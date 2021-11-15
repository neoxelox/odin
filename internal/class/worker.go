package class

import (
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/core"
)

type Worker struct {
	Configuration internal.Configuration
	Logger        core.Logger
}

func NewWorker(configuration internal.Configuration, logger core.Logger) *Worker {
	logger.SetLogger(logger.Logger().With().Str("layer", "worker").Logger())

	return &Worker{
		Configuration: configuration,
		Logger:        logger,
	}
}
