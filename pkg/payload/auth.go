package payload

import (
	"github.com/neoxelox/odin/internal/class"
	"github.com/nyaruka/phonenumbers"
)

type PostLoginStartRequest struct {
	class.Payload
	Phone string `json:"phone" validate:"required"`
}

func (self *PostLoginStartRequest) Process() error {
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

type PostLoginStartResponse struct {
	class.Payload
	ID string `json:"id"`
}

type PostLoginEndRequest struct {
	class.Payload
	ID   string `json:"id" validate:"required"`
	Code string `json:"code" validate:"required"`
}

type PostLoginEndResponse struct {
	class.Payload
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

type PostLogoutResponse struct {
	class.Payload
}
