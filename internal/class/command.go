package class

import (
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/core"
)

type Command struct {
	Configuration internal.Configuration
	Logger        core.Logger
}

func NewCommand(configuration internal.Configuration, logger core.Logger) *Command {
	logger.SetLogger(logger.Logger().With().Str("layer", "command").Logger())

	return &Command{
		Configuration: configuration,
		Logger:        logger,
	}
}
