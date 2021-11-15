package repository

import (
	"context"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
)

const USER_TABLE = "user"

var ErrUserGeneric = internal.NewError("User query failed")

type UserRepository struct {
	class.Repository
}

func NewUserRepository(configuration internal.Configuration, logger core.Logger, database database.Database) *UserRepository {
	return &UserRepository{
		Repository: *class.NewRepository(USER_TABLE, configuration, logger, database),
	}
}

func (self *UserRepository) Transaction(ctx context.Context, fn func(*UserRepository) error) error {
	return self.Database.Transaction(ctx, func(db *database.Database) error {
		return fn(&UserRepository{
			Repository: *class.NewRepository(self.Table, self.Configuration, self.Logger, *db),
		})
	})
}
