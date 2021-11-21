package repository

import (
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
)

const INVITATION_TABLE = "invitation"

var ErrInvitationGeneric = internal.NewError("Invitation query failed")

type InvitationRepository struct {
	class.Repository
}

func NewInvitationRepository(configuration internal.Configuration, logger core.Logger, database database.Database) *InvitationRepository {
	return &InvitationRepository{
		Repository: *class.NewRepository(configuration, logger, database),
	}
}
