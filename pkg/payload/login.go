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
		return ErrInvalidPhone.Cause(err)
	}

	ok := phonenumbers.IsValidNumber(ph)
	if !ok {
		return ErrInvalidPhone
	}

	self.Phone = phonenumbers.Format(ph, phonenumbers.E164)

	return nil
}

type PostLoginStartResponse struct {
	class.Payload
}

type PostLoginEndRequest struct {
	class.Payload
	Code string `json:"code"`
}

type PostLoginEndResponse struct {
	class.Payload
	AccessToken string `json:"access_token"`
	User        *User  `json:"user"`
}
