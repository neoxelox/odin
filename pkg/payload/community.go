package payload

import (
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/nyaruka/phonenumbers"
)

type Community struct {
	ID         string    `json:"id"`
	Address    string    `json:"address"`
	Name       string    `json:"name"`
	Categories []string  `json:"categories"`
	PinnedIDs  []string  `json:"pinned_ids"`
	CreatedAt  time.Time `json:"created_at"`
}

type Membership struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	CommunityID string    `json:"community_id"`
	Door        string    `json:"door"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

type CommunityAndMembership struct {
	Community  Community  `json:"community"`
	Membership Membership `json:"membership"`
}

type UserAndMembership struct {
	User       User       `json:"user"`
	Membership Membership `json:"membership"`
}

type PostCommunityRequest struct {
	class.Payload
	Address    string    `json:"address" validate:"required"`
	Name       *string   `json:"name" validate:"omitempty,required"`
	Categories *[]string `json:"categories" validate:"omitempty,required"`
}

type PostCommunityResponse struct {
	class.Payload
	CommunityAndMembership
}

type GetCommunityRequest struct {
	class.Payload
	ID string `param:"id" validate:"required"`
}

type GetCommunityResponse struct {
	class.Payload
	CommunityAndMembership
}

type GetCommunityListResponse struct {
	class.Payload
	Communities []CommunityAndMembership `json:"communities"`
}

type GetCommunityUserRequest struct {
	class.Payload
	CommunityID  string `param:"community_id" validate:"required"`
	MembershipID string `param:"membership_id" validate:"required"`
}

type GetCommunityUserResponse struct {
	class.Payload
	UserAndMembership
}

type GetCommunityUserListRequest struct {
	class.Payload
	ID string `param:"id" validate:"required"`
}

type GetCommunityUserListResponse struct {
	class.Payload
	Users []UserAndMembership `json:"users"`
}

type PostCommunityInviteRequest struct {
	class.Payload
	ID    string `param:"id" validate:"required"`
	Phone string `json:"phone" validate:"required"`
	Door  string `json:"door"`
	Role  string `json:"role" validate:"required"`
}

func (self *PostCommunityInviteRequest) Process() error {
	ph, err := phonenumbers.Parse(self.Phone, "ES")
	if err != nil {
		return ExcInvalidPhone.Cause(err)
	}

	if !phonenumbers.IsValidNumber(ph) {
		return ExcInvalidPhone
	}

	self.Phone = phonenumbers.Format(ph, phonenumbers.E164)

	return nil
}

type PostCommunityInviteResponse struct {
	class.Payload
	Invitation
}

type PostCommunityLeaveRequest struct {
	class.Payload
	ID string `param:"id" validate:"required"`
}

type PostCommunityLeaveResponse struct {
	class.Payload
}
