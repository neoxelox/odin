package model

import (
	"fmt"
	"time"

	"github.com/aodin/date"

	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/utility"
	"github.com/rs/xid"
)

const (
	USER_NAME_MAX_LENGTH = 100
	USER_NAME_MIN_LENGTH = 1
	USER_DEFAULT_PICTURE = "https://eu.ui-avatars.com/api/?name=%s&size=128"
)

type User struct {
	class.Model
	ID            string     `db:"id"`
	Phone         string     `db:"phone"`
	Name          string     `db:"name"`
	Email         string     `db:"email"`
	Picture       string     `db:"picture"`
	Birthday      date.Date  `db:"birthday"`
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

func (self *User) Copy() *User {
	return &User{
		ID:            *utility.CopyString(&self.ID),
		Phone:         *utility.CopyString(&self.Phone),
		Name:          *utility.CopyString(&self.Name),
		Email:         *utility.CopyString(&self.Email),
		Picture:       *utility.CopyString(&self.Picture),
		Birthday:      *utility.CopyDate(&self.Birthday),
		Language:      *utility.CopyString(&self.Language),
		LastSessionID: utility.CopyString(self.LastSessionID),
		IsBanned:      *utility.CopyBool(&self.IsBanned),
		CreatedAt:     *utility.CopyTime(&self.CreatedAt),
		DeletedAt:     utility.CopyTime(self.DeletedAt),
	}
}
