package repository

import (
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
)

const MEMBERSHIP_TABLE = "membership"

var ErrMembershipGeneric = internal.NewError("Membership query failed")

type MembershipRepository struct {
	class.Repository
}

func NewMembershipRepository(configuration internal.Configuration, logger core.Logger, database database.Database) *MembershipRepository {
	return &MembershipRepository{
		Repository: *class.NewRepository(configuration, logger, database),
	}
}
