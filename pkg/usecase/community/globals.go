package community

import "github.com/neoxelox/odin/internal"

var (
	ErrGeneric        = internal.NewError("Community execution failed")
	ErrInvalid        = internal.NewError("Community is invalid")
	ErrInvalidAddress = internal.NewError("Community address invalid")
	ErrInvalidName    = internal.NewError("Community name invalid")
	ErrInvalidDoor    = internal.NewError("Membership door invalid")
	ErrInvalidRole    = internal.NewError("Membership role invalid")
	ErrAlreadyJoined  = internal.NewError("User already joined this community")
	ErrNotBelongs     = internal.NewError("User does not belong to this community")
	ErrNotPermission  = internal.NewError("User does not have permission to perform this action")
)
