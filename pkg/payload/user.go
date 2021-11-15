package payload

import (
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/rs/xid"
)

type User struct {
	ID       xid.ID     `json:"id"`
	Phone    string     `json:"phone"`
	Name     string     `json:"name"`
	Email    *string    `json:"email"`
	Picture  *string    `json:"picture"`
	Birthday *time.Time `json:"birthday"`
}

type GetUserRequest struct {
	class.Payload
	ID xid.ID `param:"id"`
}

type GetUserResponse struct {
	class.Payload
	User
}

type PostUserProfileRequest struct {
	class.Payload
	Name     string     `json:"name"`
	LastName string     `json:"last_name"`
	Picture  *string    `json:"picture"`
	Birthday *time.Time `json:"birthday"`
}

type PostUserProfileResponse struct {
	class.Payload
	Name     string     `json:"name"`
	LastName string     `json:"last_name"`
	Picture  *string    `json:"picture"`
	Birthday *time.Time `json:"birthday"`
}

type PostUserEmailStartRequest struct {
	class.Payload
	Email string `json:"email"`
}

type PostUserEmailStartResponse struct {
	class.Payload
}

type PostUserEmailEndRequest struct {
	class.Payload
	Code string `json:"code"`
}

type PostUserEmailEndResponse struct {
	class.Payload
	Email string `json:"email"`
}

type PostUserPhoneStartRequest struct {
	class.Payload
	Phone string `json:"phone"`
}

type PostUserPhoneStartResponse struct {
	class.Payload
}

type PostUserPhoneEndRequest struct {
	class.Payload
	Code string `json:"code"`
}

type PostUserPhoneEndResponse struct {
	class.Payload
	Phone string `json:"phone"`
}
