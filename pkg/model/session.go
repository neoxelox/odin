package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/rs/xid"
)

type Session struct {
	class.Model
	ID         string          `db:"id"`
	UserID     string          `db:"user_id"`
	Metadata   SessionMetadata `db:"metadata"`
	CreatedAt  time.Time       `db:"created_at"`
	LastSeenAt time.Time       `db:"last_seen_at"`
}

type SessionMetadata struct {
	IP         string
	Device     string
	ApiVersion string
}

func NewSession() *Session {
	now := time.Now()

	return &Session{
		ID:         xid.New().String(),
		CreatedAt:  now,
		LastSeenAt: now,
	}
}

func (self Session) String() string {
	return fmt.Sprintf("<%s: %s>", self.UserID, self.ID)
}
