package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/rs/xid"
)

type Community struct {
	class.Model
	Address    string   `db:"address"`
	Name       string   `db:"name"`
	Categories []string `db:"categories"`
	PinnedIDs  []xid.ID `db:"pinned_ids"`
}

func NewCommunity() *Community {
	now := time.Now()

	return &Community{
		Model: class.Model{
			ID:        xid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (self Community) String() string {
	return fmt.Sprintf("<%s: %s>", self.Name, self.ID)
}
