package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/rs/xid"
)

type Membership struct {
	class.Model
	Phone       *string `db:"phone"`
	UserID      *xid.ID `db:"user_id"`
	CommunityID xid.ID  `db:"community_id"`
	State       string  `db:"state"`
	Door        string  `db:"door"`
	Role        string  `db:"role"`
}

func NewMembership() *Membership {
	now := time.Now()

	return &Membership{
		Model: class.Model{
			ID:        xid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (self Membership) String() string {
	str := "Nil"
	if self.UserID != nil {
		str = self.UserID.String()
	} else if self.Phone != nil {
		str = *self.Phone
	}

	return fmt.Sprintf("<%s <-> %s: %s>", str, self.CommunityID, self.ID)
}
