package internal

import (
	"encoding/json"

	"github.com/cockroachdb/errors"
)

type Exception struct {
	Origin  error  `json:"-"`
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

func NewException(status int, code string) *Exception {
	return &Exception{
		Status: status,
		Code:   code,
	}
}

func (self *Exception) Cause(err error) *Exception {
	self.Origin = errors.WrapWithDepth(1, err, self.Code)
	self.Message = err.Error()
	return self
}

func (self *Exception) Redact() {
	self.Message = ""
}

func (self Exception) Error() string {
	return self.Origin.Error()
}

func (self Exception) String() string {
	return self.Code
}

func (self Exception) Unwrap() error {
	return self.Origin
}

func (self Exception) JSON() string {
	json, _ := json.Marshal(self)
	return string(json)
}
