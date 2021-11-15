package repository

import (
	"context"

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
		Repository: *class.NewRepository(COMMUNITY_TABLE, configuration, logger, database),
	}
}

func (self *CommunityRepository) Transaction(ctx context.Context, fn func(*CommunityRepository) error) error {
	return self.Database.Transaction(ctx, func(db *database.Database) error {
		return fn(&CommunityRepository{
			Repository: *class.NewRepository(self.Table, self.Configuration, self.Logger, *db),
		})
	})
}
