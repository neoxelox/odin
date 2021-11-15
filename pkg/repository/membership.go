package repository

import (
	"context"

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
		Repository: *class.NewRepository(MEMBERSHIP_TABLE, configuration, logger, database),
	}
}

func (self *MembershipRepository) Transaction(ctx context.Context, fn func(*MembershipRepository) error) error {
	return self.Database.Transaction(ctx, func(db *database.Database) error {
		return fn(&MembershipRepository{
			Repository: *class.NewRepository(self.Table, self.Configuration, self.Logger, *db),
		})
	})
}
