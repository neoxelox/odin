package repository

import (
	"context"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
)

const POST_HISTORY_TABLE = "post_history"

var ErrPostHistoryGeneric = internal.NewError("Post History query failed")

type PostHistoryRepository struct {
	class.Repository
}

func NewPostHistoryRepository(configuration internal.Configuration, logger core.Logger, database database.Database) *PostHistoryRepository {
	return &PostHistoryRepository{
		Repository: *class.NewRepository(POST_HISTORY_TABLE, configuration, logger, database),
	}
}

func (self *PostHistoryRepository) Transaction(ctx context.Context, fn func(*PostHistoryRepository) error) error {
	return self.Database.Transaction(ctx, func(db *database.Database) error {
		return fn(&PostHistoryRepository{
			Repository: *class.NewRepository(self.Table, self.Configuration, self.Logger, *db),
		})
	})
}
