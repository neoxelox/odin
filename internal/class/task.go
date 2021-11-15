package class

import (
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/core"
)

type Task struct {
	Configuration internal.Configuration
	Logger        core.Logger
}

func NewTask(configuration internal.Configuration, logger core.Logger) *Task {
	logger.SetLogger(logger.Logger().With().Str("layer", "task").Logger())

	return &Task{
		Configuration: configuration,
		Logger:        logger,
	}
}
