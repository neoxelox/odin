package repository

import (
	"context"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
)

const SESSION_TABLE = "session"

var ErrSessionGeneric = internal.NewError("Session query failed")

type SessionRepository struct {
	class.Repository
}

func NewSessionRepository(configuration internal.Configuration, logger core.Logger, database database.Database) *SessionRepository {
	return &SessionRepository{
		Repository: *class.NewRepository(SESSION_TABLE, configuration, logger, database),
	}
}

func (self *SessionRepository) Transaction(ctx context.Context, fn func(*SessionRepository) error) error {
	return self.Database.Transaction(ctx, func(db *database.Database) error {
		return fn(&SessionRepository{
			Repository: *class.NewRepository(self.Table, self.Configuration, self.Logger, *db),
		})
	})
}
