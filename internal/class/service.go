package class

import (
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/core"
)

type Service struct {
	Configuration internal.Configuration
	Logger        core.Logger
}

func NewService(configuration internal.Configuration, logger core.Logger) *Service {
	logger.SetLogger(logger.Logger().With().Str("layer", "service").Logger())

	return &Service{
		Configuration: configuration,
		Logger:        logger,
	}
}
