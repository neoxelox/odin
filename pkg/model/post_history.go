package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/rs/xid"
)

type PostHistory struct {
	class.Model
	PostID     xid.ID             `db:"post_id"`
	Message    string             `db:"message"`
	Categories []string           `db:"categories"`
	State      *string            `db:"state"`
	Widgets    PostHistoryWidgets `db:"widgets"`
	Media      PostHistoryMedia   `db:"media"`
}

type PostHistoryWidgets struct {
	Poll *map[string][]xid.ID
}

type PostHistoryMedia struct {
	Pictures *[]string
	Videos   *[]string
	Audios   *[]string
}

func NewPostHistory() *PostHistory {
	now := time.Now()

	return &PostHistory{
		Model: class.Model{
			ID:        xid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (self PostHistory) String() string {
	return fmt.Sprintf("<%s: %s>", self.Message, self.ID)
}
