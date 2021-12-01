package payload

import (
	"time"

	"github.com/neoxelox/odin/internal/class"
)

type Invitation struct {
	ID          string    `json:"id"`
	Phone       string    `json:"phone"`
	CommunityID string    `json:"community_id"`
	Door        string    `json:"door"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetInvitationListResponse struct {
	class.Payload
	Invitations []Invitation `json:"invitations"`
}

type PostInvitationAcceptRequest struct {
	class.Payload
	ID string `param:"id" validate:"required"`
}

type PostInvitationAcceptResponse struct {
	class.Payload
	Membership
}

type PostInvitationRejectRequest struct {
	class.Payload
	ID string `param:"id" validate:"required"`
}

type PostInvitationRejectResponse struct {
	class.Payload
}
