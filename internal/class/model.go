package class

import (
	"time"

	"github.com/rs/xid"
)

type Model struct {
	ID        xid.ID     `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
