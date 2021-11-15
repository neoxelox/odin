package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/rs/xid"
)

type OTP struct {
	class.Model
	Asset    string `db:"asset"`
	Type     string `db:"type"`
	Code     string `db:"code"`
	Attempts int    `db:"attempts"`
}

var OTPType = struct {
	PHONE string
	EMAIL string
}{"PHONE", "EMAIL"}

func NewOTP() *OTP {
	now := time.Now()

	return &OTP{
		Model: class.Model{
			ID:        xid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (self OTP) String() string {
	return fmt.Sprintf("<%s: %s>", self.Asset, self.ID)
}
