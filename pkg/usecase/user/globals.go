package user

import "github.com/neoxelox/odin/internal"

var (
	ErrGeneric      = internal.NewError("User execution failed")
	ErrInvalidPhone = internal.NewError("User phone invalid")
	ErrPhoneExists  = internal.NewError("Another user with the same phone already exists")
)
