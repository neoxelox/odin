package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/utility"
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

func (self *Membership) Copy() *Membership {
	return &Membership{
		ID:          *utility.CopyString(&self.ID),
		UserID:      *utility.CopyString(&self.UserID),
		CommunityID: *utility.CopyString(&self.CommunityID),
		Door:        *utility.CopyString(&self.Door),
		Role:        *utility.CopyString(&self.Role),
		CreatedAt:   *utility.CopyTime(&self.CreatedAt),
		DeletedAt:   utility.CopyTime(self.DeletedAt),
	}
}
