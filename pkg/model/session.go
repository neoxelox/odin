package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/rs/xid"
)

type Session struct {
	class.Model
	UserID   xid.ID          `db:"user_id"`
	Metadata SessionMetadata `db:"metadata"`
}

type SessionMetadata struct {
	IP     string
	Device string
}

func NewSession() *Session {
	now := time.Now()

	return &Session{
		Model: class.Model{
			ID:        xid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (self Session) String() string {
	return fmt.Sprintf("<%s: %s>", self.UserID, self.ID)
}
