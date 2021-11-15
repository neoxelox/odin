package class

import (
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
)

type Repository struct {
	Table         string
	Configuration internal.Configuration
	Logger        core.Logger
	Database      database.Database
}

func NewRepository(table string, configuration internal.Configuration, logger core.Logger, database database.Database) *Repository {
	logger.SetLogger(logger.Logger().With().Str("layer", "repository").Logger())

	return &Repository{
		Table:         table,
		Configuration: configuration,
		Logger:        logger,
		Database:      database,
	}
}
