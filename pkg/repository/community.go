package repository

import (
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
)

const COMMUNITY_TABLE = "community"

var ErrCommunityGeneric = internal.NewError("Community query failed")

type CommunityRepository struct {
	class.Repository
}

func NewCommunityRepository(configuration internal.Configuration, logger core.Logger, database database.Database) *CommunityRepository {
	return &CommunityRepository{
		Repository: *class.NewRepository(configuration, logger, database),
	}
}
