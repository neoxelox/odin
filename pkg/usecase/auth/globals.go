package auth

import (
	"regexp"

	"github.com/neoxelox/odin/internal"
)

var (
	ErrGeneric             = internal.NewError("Auth execution failed")
	ErrTamperedAccessToken = internal.NewError("Access token decryption failed")
	ErrExpiredAccessToken  = internal.NewError("Access token is expired")
	ErrInvalidAccessToken  = internal.NewError("Access token is invalid")
	ErrExpiredSession      = internal.NewError("Session has expired")
	ErrDeletedUser         = internal.NewError("User is deleted")
	ErrBannedUser          = internal.NewError("User is banned")
)

var UNVERSIONED_PATHS = regexp.MustCompile(`^/(login|logout|file).*$`)
