package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/rs/xid"
)

type User struct {
	class.Model
	Phone         string     `db:"phone"`
	Name          string     `db:"name"`
	Email         *string    `db:"email"`
	Picture       *string    `db:"picture"`
	Birthday      *time.Time `db:"birthday"`
	Language      string     `db:"language"`
	LastSessionID *xid.ID    `db:"last_session_id"`
	IsBanned      bool       `db:"is_banned"`
}

func NewUser() *User {
	now := time.Now()

	return &User{
		Model: class.Model{
			ID:        xid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Language: "ES",
	}
}

func (self User) String() string {
	return fmt.Sprintf("<%s: %s>", self.Name, self.ID)
}
