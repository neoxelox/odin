package auth

import "github.com/neoxelox/odin/internal"

var (
	ErrGeneric             = internal.NewError("Auth execution failed")
	ErrTamperedAccessToken = internal.NewError("Access token decryption failed")
	ErrExpiredAccessToken  = internal.NewError("Access token is expired")
	ErrInvalidAccessToken  = internal.NewError("Access token is invalid")
	ErrUserBanned          = internal.NewError("User is banned")
)
