package repository

import (
	"context"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
)

const POST_TABLE = "post"

var ErrPostGeneric = internal.NewError("Post query failed")

type PostRepository struct {
	class.Repository
}

func NewPostRepository(configuration internal.Configuration, logger core.Logger, database database.Database) *PostRepository {
	return &PostRepository{
		Repository: *class.NewRepository(POST_TABLE, configuration, logger, database),
	}
}

func (self *PostRepository) Transaction(ctx context.Context, fn func(*PostRepository) error) error {
	return self.Database.Transaction(ctx, func(db *database.Database) error {
		return fn(&PostRepository{
			Repository: *class.NewRepository(self.Table, self.Configuration, self.Logger, *db),
		})
	})
}
