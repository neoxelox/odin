package auth

import "github.com/neoxelox/odin/internal"

var (
	ErrGeneric             = internal.NewError("Auth execution failed")
	ErrTamperedAccessToken = internal.NewError("Access token decryption failed")
	ErrExpiredAccessToken  = internal.NewError("Access token is expired")
	ErrInvalidAccessToken  = internal.NewError("Access token is invalid")
	ErrExpiredSession      = internal.NewError("Session has expired")
	ErrBannedUser          = internal.NewError("User is banned")
)

var UNVERSIONED_PATHS = []string{
	"/login/start",
	"/login/end",
	"/logout",
}
