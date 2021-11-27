package payload

import (
	"github.com/aodin/date"
	"github.com/badoux/checkmail"
	"github.com/neoxelox/odin/internal/class"
	"github.com/nyaruka/phonenumbers"
)

type User struct {
	ID       string    `json:"id"`
	Phone    string    `json:"phone"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Picture  string    `json:"picture"`
	Birthday date.Date `json:"birthday"`
}

type GetUserProfileResponse struct {
	class.Payload
	User
}

type PostUserProfileRequest struct {
	class.Payload
	Name     *string    `json:"name" validate:"omitempty,required"`
	LastName *string    `json:"last_name" validate:"omitempty,required"`
	Picture  *string    `json:"picture" validate:"omitempty,required"`
	Birthday *date.Date `json:"birthday" validate:"omitempty,required"`
}

type PostUserProfileResponse struct {
	class.Payload
	User
}

type PostUserEmailStartRequest struct {
	class.Payload
	Email string `json:"email" validate:"required"`
}

func (self *PostUserEmailStartRequest) Process() error {
	err := checkmail.ValidateFormat(self.Email)
	if err != nil {
		return ExcInvalidEmail.Cause(err)
	}

	return nil
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

type PostUserEmailEndResponse struct {
	class.Payload
	Email string `json:"email"`
}

type PostUserPhoneStartRequest struct {
	class.Payload
	Phone string `json:"phone" validate:"required"`
}

func (self *PostUserPhoneStartRequest) Process() error {
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

type PostUserPhoneStartResponse struct {
	class.Payload
	ID string `json:"id"`
}

type PostUserPhoneEndRequest struct {
	class.Payload
	ID   string `json:"id" validate:"required"`
	Code string `json:"code" validate:"required"`
}

type PostUserPhoneEndResponse struct {
	class.Payload
	Phone string `json:"phone"`
}

type DeleteUserResponse struct {
	class.Payload
}
