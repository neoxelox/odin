package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/rs/xid"
)

type Community struct {
	class.Model
	ID         string     `db:"id"`
	Address    string     `db:"address"`
	Name       string     `db:"name"`
	Categories []string   `db:"categories"`
	PinnedIDs  []string   `db:"pinned_ids"`
	CreatedAt  time.Time  `db:"created_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}

func NewCommunity() *Community {
	now := time.Now()

	return &Community{
		ID:        xid.New().String(),
		CreatedAt: now,
	}
}

func (self Community) String() string {
	return fmt.Sprintf("<%s: %s>", self.Name, self.ID)
}
