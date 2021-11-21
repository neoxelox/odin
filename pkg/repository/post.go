package repository

import (
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
)

const (
	POST_TABLE         = "post"
	POST_HISTORY_TABLE = "post_history"
)

var ErrPostGeneric = internal.NewError("Post query failed")

type PostRepository struct {
	class.Repository
}

func NewPostRepository(configuration internal.Configuration, logger core.Logger, database database.Database) *PostRepository {
	return &PostRepository{
		Repository: *class.NewRepository(configuration, logger, database),
	}
}
