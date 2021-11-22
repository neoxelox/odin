package user

import "github.com/neoxelox/odin/internal"

var (
	ErrGeneric         = internal.NewError("User execution failed")
	ErrInvalidPhone    = internal.NewError("User phone invalid")
	ErrInvalidName     = internal.NewError("User name invalid")
	ErrInvalidPicture  = internal.NewError("User picture invalid")
	ErrInvalidBirthday = internal.NewError("User birthday invalid")
	ErrInvalidEmail    = internal.NewError("User email invalid")
	ErrPhoneExists     = internal.NewError("Another user with the same phone already exists")
)
