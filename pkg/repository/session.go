package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/pkg/model"
)

const SESSION_TABLE = "session"

var ErrSessionGeneric = internal.NewError("Session query failed")

type SessionRepository struct {
	class.Repository
}

func NewSessionRepository(configuration internal.Configuration, logger core.Logger, database database.Database) *SessionRepository {
	return &SessionRepository{
		Repository: *class.NewRepository(configuration, logger, database),
	}
}

func (self *SessionRepository) Create(ctx context.Context, session model.Session) (*model.Session, error) {
	var s model.Session

	query := fmt.Sprintf(`INSERT INTO "%s"
						  ("id", "user_id", "metadata", "created_at", "last_seen_at", "expired_at")
						  VALUES ($1, $2, $3, $4, $5, $6)
						  RETURNING *;`, SESSION_TABLE)

	err := self.Database.Query(
		ctx, query, session.ID, session.UserID, session.Metadata, session.CreatedAt, session.LastSeenAt, session.ExpiredAt).Scan(&s)
	if err != nil {
		return nil, ErrSessionGeneric().Wrap(err)
	}

	return &s, nil
}

func (self *SessionRepository) GetByID(ctx context.Context, id string) (*model.Session, error) {
	var s model.Session

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "id" = $1;`, SESSION_TABLE)

	err := self.Database.Query(ctx, query, id).Scan(&s)
	switch {
	case err == nil:
		return &s, nil
	case database.ErrNoRows().Is(err):
		return nil, nil
	default:
		return nil, ErrSessionGeneric().Wrap(err)
	}
}

func (self *SessionRepository) UpdateLastSeen(ctx context.Context, id string, lastSeen time.Time) error {
	query := fmt.Sprintf(`UPDATE "%s"
						  SET "last_seen_at" = $1
						  WHERE "id" = $2;`, SESSION_TABLE)

	affected, err := self.Database.Exec(ctx, query, lastSeen, id)
	if err != nil {
		return ErrSessionGeneric().Wrap(err)
	}

	if affected != 1 {
		return ErrSessionGeneric()
	}

	return nil
}

func (self *SessionRepository) UpdateExpiredAt(ctx context.Context, id string, expiredAt *time.Time) error {
	query := fmt.Sprintf(`UPDATE "%s"
						  SET "expired_at" = $1
						  WHERE "id" = $2;`, SESSION_TABLE)

	affected, err := self.Database.Exec(ctx, query, expiredAt, id)
	if err != nil {
		return ErrSessionGeneric().Wrap(err)
	}

	if affected != 1 {
		return ErrSessionGeneric()
	}

	return nil
}
