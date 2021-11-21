package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/rs/xid"
)

type User struct {
	class.Model
	ID            string     `db:"id"`
	Phone         string     `db:"phone"`
	Name          string     `db:"name"`
	Email         *string    `db:"email"`
	Picture       *string    `db:"picture"`
	Birthday      *time.Time `db:"birthday"`
	Language      string     `db:"language"`
	LastSessionID *string    `db:"last_session_id"`
	IsBanned      bool       `db:"is_banned"`
	CreatedAt     time.Time  `db:"created_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}

func NewUser() *User {
	now := time.Now()

	return &User{
		ID:        xid.New().String(),
		CreatedAt: now,
		Language:  "ES",
	}
}

func (self User) String() string {
	return fmt.Sprintf("<%s: %s>", self.Name, self.ID)
}
