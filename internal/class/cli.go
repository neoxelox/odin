package class

import (
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/core"
)

type CLI struct {
	Configuration internal.Configuration
	Logger        core.Logger
}

func NewCLI(configuration internal.Configuration, logger core.Logger) *CLI {
	logger.SetLogger(logger.Logger().With().Str("layer", "cli").Logger())

	return &CLI{
		Configuration: configuration,
		Logger:        logger,
	}
}
