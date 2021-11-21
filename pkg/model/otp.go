package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/rs/xid"
)

const (
	OTP_CODE_LENGTH  = 6
	OTP_MAX_ATTEMPTS = 5
	OTP_EXPIRATION   = time.Duration(5) * time.Minute
)

type OTP struct {
	class.Model
	ID        string    `db:"id"`
	Asset     string    `db:"asset"`
	Type      string    `db:"type"`
	Code      string    `db:"code"`
	Attempts  int       `db:"attempts"`
	ExpiresAt time.Time `db:"expires_at"`
}

var OTPType = struct {
	SMS   string
	EMAIL string
	Has   func(typee string) bool
}{"SMS", "EMAIL", func(typee string) bool {
	return typee == "SMS" || typee == "EMAIL"
}}

func NewOTP() *OTP {
	now := time.Now()

	return &OTP{
		ID:        xid.New().String(),
		ExpiresAt: now.Add(OTP_EXPIRATION),
	}
}

func (self OTP) String() string {
	return fmt.Sprintf("<%s: %s>", self.Asset, self.ID)
}
