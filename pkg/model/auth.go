package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/utility"
)

const (
	ACCESS_TOKEN_EXPIRATION                     = time.Duration(24*365) * time.Hour
	CONTEXT_SESSION_KEY     internal.ContextKey = "auth:session"
	CONTEXT_USER_KEY        internal.ContextKey = "auth:user"
)

type AccessToken struct {
	class.Model
	Private AccessTokenPrivate
	Public  AccessTokenPublic
}

type AccessTokenPrivate struct {
	SessionID string    `json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
	// Redundant field in order not to hit the DB to validate session expiration
	ExpiresAt time.Time `json:"expires_at"`
}

// Needed for Paseto library (Validation will be performed in the usecase)
func (self *AccessTokenPrivate) Valid() error { return nil }

type AccessTokenPublic struct {
	ApiVersion string `json:"api_version"`
}

func NewAccessToken() *AccessToken {
	now := time.Now()

	return &AccessToken{
		Private: AccessTokenPrivate{
			CreatedAt: now,
			ExpiresAt: now.Add(ACCESS_TOKEN_EXPIRATION),
		},
		Public: AccessTokenPublic{},
	}
}

func (self AccessToken) String() string {
	return fmt.Sprintf("<%s: %s>", self.Public.ApiVersion, self.Private.SessionID)
}

func (self *AccessToken) Copy() *AccessToken {
	return &AccessToken{
		Private: AccessTokenPrivate{
			SessionID: *utility.CopyString(&self.Private.SessionID),
			CreatedAt: *utility.CopyTime(&self.Private.CreatedAt),
			ExpiresAt: *utility.CopyTime(&self.Private.ExpiresAt),
		},
		Public: AccessTokenPublic{
			ApiVersion: *utility.CopyString(&self.Public.ApiVersion),
		},
	}
}
