package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/rs/xid"
)

type Post struct {
	class.Model
	ThreadID      *xid.ID   `db:"thread_id"`
	CreatorID     xid.ID    `db:"creator_id"`
	LastHistoryID *xid.ID   `db:"last_history_id"`
	Type          string    `db:"type"`
	Priority      *int      `db:"priority"`
	RecipientIDs  *[]xid.ID `db:"recipient_ids"`
	VoterIDs      []xid.ID  `db:"voter_ids"`
}

func NewPost() *Post {
	now := time.Now()

	return &Post{
		Model: class.Model{
			ID:        xid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (self Post) String() string {
	return fmt.Sprintf("<%s: %s>", self.Type, self.ID)
}
