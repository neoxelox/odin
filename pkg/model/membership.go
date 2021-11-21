package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/rs/xid"
)

type Membership struct {
	class.Model
	ID          string     `db:"id"`
	UserID      string     `db:"user_id"`
	CommunityID string     `db:"community_id"`
	Door        string     `db:"door"`
	Role        string     `db:"role"`
	CreatedAt   time.Time  `db:"created_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

var MembershipRole = struct {
	ADMINISTRATOR string
	PRESIDENT     string
	SECRETARY     string
	RESIDENT      string
	LESSEE        string
	Has           func(role string) bool
}{"ADMINISTRATOR", "PRESIDENT", "SECRETARY", "RESIDENT", "LESSEE",
	func(role string) bool {
		return role == "ADMINISTRATOR" || role == "PRESIDENT" ||
			role == "SECRETARY" || role == "RESIDENT" || role == "LESSEE"
	},
}

func NewMembership() *Membership {
	now := time.Now()

	return &Membership{
		ID:        xid.New().String(),
		CreatedAt: now,
	}
}

func (self Membership) String() string {
	return fmt.Sprintf("<%s <-> %s: %s>", self.UserID, self.CommunityID, self.ID)
}
