package invitation

import "github.com/neoxelox/odin/internal"

var (
	ErrGeneric          = internal.NewError("Invitation execution failed")
	ErrInvalid          = internal.NewError("Invitation is invalid")
	ErrExpired          = internal.NewError("Invitation has expired")
	ErrInvitingYourself = internal.NewError("User is inviting itself")
	ErrAlreadyInvited   = internal.NewError("User has already been invited")
)
