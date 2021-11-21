package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/rs/xid"
)

const (
	INVITATION_EXPIRATION = time.Duration(24*30) * time.Hour
)

type Invitation struct {
	class.Model
	ID          string    `db:"id"`
	Phone       string    `db:"phone"`
	CommunityID string    `db:"community_id"`
	State       string    `db:"state"`
	Door        string    `db:"door"`
	Role        string    `db:"role"`
	CreatedAt   time.Time `db:"created_at"`
	RemindedAt  time.Time `db:"reminded_at"`
	ExpiresAt   time.Time `db:"expires_at"`
}

var InvitationState = struct {
	PENDING  string
	ACCEPTED string
	REJECTED string
	Has      func(state string) bool
}{"PENDING", "ACCEPTED", "REJECTED", func(state string) bool {
	return state == "PENDING" || state == "ACCEPTED" || state == "REJECTED"
}}

func NewInvitation() *Invitation {
	now := time.Now()

	return &Invitation{
		ID:         xid.New().String(),
		CreatedAt:  now,
		RemindedAt: now,
	}
}

func (self Invitation) String() string {
	return fmt.Sprintf("<%s <-> %s: %s>", self.Phone, self.CommunityID, self.ID)
}
