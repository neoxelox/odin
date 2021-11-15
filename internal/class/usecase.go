package class

import (
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/core"
)

type Usecase struct {
	Configuration internal.Configuration
	Logger        core.Logger
}

func NewUsecase(configuration internal.Configuration, logger core.Logger) *Usecase {
	logger.SetLogger(logger.Logger().With().Str("layer", "usecase").Logger())

	return &Usecase{
		Configuration: configuration,
		Logger:        logger,
	}
}
