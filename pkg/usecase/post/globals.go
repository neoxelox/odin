package post

import "github.com/neoxelox/odin/internal"

var (
	ErrGeneric           = internal.NewError("Post execution failed")
	ErrInvalid           = internal.NewError("Post is invalid")
	ErrInvalidType       = internal.NewError("Post type invalid")
	ErrInvalidThread     = internal.NewError("Publication thread invalid")
	ErrInvalidPriority   = internal.NewError("Issue priority invalid")
	ErrInvalidRecipients = internal.NewError("Post recipients invalid")
	ErrInvalidMessage    = internal.NewError("Post message invalid")
	ErrInvalidState      = internal.NewError("Issue state invalid")
	ErrInvalidMedia      = internal.NewError("Post media invalid")
	ErrInvalidPoll       = internal.NewError("Publication poll widget invalid")
)
