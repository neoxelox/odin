package payload

import (
	"time"

	"github.com/neoxelox/odin/internal/class"
)

type User struct {
	ID       string     `json:"id"`
	Phone    string     `json:"phone"`
	Name     string     `json:"name"`
	Email    *string    `json:"email"`
	Picture  *string    `json:"picture"`
	Birthday *time.Time `json:"birthday"`
}

type GetUserRequest struct {
	class.Payload
	ID string `param:"id" validate:"required"`
}

type GetUserResponse struct {
	class.Payload
	User
}

type PostUserProfileRequest struct {
	class.Payload
	Name     string     `json:"name" validate:"required"`
	LastName string     `json:"last_name" validate:"required"`
	Picture  *string    `json:"picture" validate:"required"`
	Birthday *time.Time `json:"birthday" validate:"required"`
}

type PostUserEmailStartRequest struct {
	class.Payload
	Email string `json:"email" validate:"required"`
}

type PostUserEmailStartResponse struct {
	class.Payload
	ID string `json:"id"`
}

type PostUserEmailEndRequest struct {
	class.Payload
	ID   string `json:"id" validate:"required"`
	Code string `json:"code" validate:"required"`
}

type PostUserPhoneStartRequest struct {
	class.Payload
	Phone string `json:"phone" validate:"required"`
}

type PostUserPhoneStartResponse struct {
	class.Payload
	ID string `json:"id"`
}

type PostUserPhoneEndRequest struct {
	class.Payload
	ID   string `json:"id" validate:"required"`
	Code string `json:"code" validate:"required"`
}
