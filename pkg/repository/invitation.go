package repository

import (
	"context"
	"fmt"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/pkg/model"
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

func (self *InvitationRepository) Create(ctx context.Context, invitation model.Invitation) (*model.Invitation, error) {
	var i model.Invitation

	query := fmt.Sprintf(`INSERT INTO "%s"
						  ("id", "phone", "community_id", "door", "role", "created_at", "reminded_at", "expires_at")
						  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
						  RETURNING *;`, INVITATION_TABLE)

	err := self.Database.Query(
		ctx, query, invitation.ID, invitation.Phone, invitation.CommunityID, invitation.Door, invitation.Role, invitation.CreatedAt, invitation.RemindedAt, invitation.ExpiresAt).Scan(&i)
	if err != nil {
		return nil, ErrInvitationGeneric().Wrap(err)
	}

	return &i, nil
}

func (self *InvitationRepository) GetByID(ctx context.Context, id string) (*model.Invitation, error) {
	var i model.Invitation

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "id" = $1;`, INVITATION_TABLE)

	err := self.Database.Query(ctx, query, id).Scan(&i)
	switch {
	case err == nil:
		return &i, nil
	case database.ErrNoRows().Is(err):
		return nil, nil
	default:
		return nil, ErrInvitationGeneric().Wrap(err)
	}
}

func (self *InvitationRepository) GetByPhone(ctx context.Context, phone string) (*model.Invitation, error) {
	var i model.Invitation

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "phone" = $1;`, INVITATION_TABLE)

	err := self.Database.Query(ctx, query, phone).Scan(&i)
	switch {
	case err == nil:
		return &i, nil
	case database.ErrNoRows().Is(err):
		return nil, nil
	default:
		return nil, ErrInvitationGeneric().Wrap(err)
	}
}

func (self *InvitationRepository) List(ctx context.Context, phone string) ([]model.Invitation, error) {
	var is []model.Invitation

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "phone" = $1;`, INVITATION_TABLE)

	err := self.Database.Query(ctx, query, phone).Scan(&is)
	switch {
	case err == nil:
		return is, nil
	case database.ErrNoRows().Is(err):
		return []model.Invitation{}, nil
	default:
		return nil, ErrInvitationGeneric().Wrap(err)
	}
}

func (self *InvitationRepository) DeleteByID(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM "%s"
						  WHERE "id" = $1;`, INVITATION_TABLE)

	affected, err := self.Database.Exec(ctx, query, id)
	if err != nil {
		return ErrInvitationGeneric().Wrap(err)
	}

	if affected != 1 {
		return ErrInvitationGeneric()
	}

	return nil
}

func (self *InvitationRepository) DeleteByIDs(ctx context.Context, ids []string) error {
	query := fmt.Sprintf(`DELETE FROM "%s"
						  WHERE "id" = ANY ($1);`, INVITATION_TABLE)

	affected, err := self.Database.Exec(ctx, query, ids)
	if err != nil {
		return ErrInvitationGeneric().Wrap(err)
	}

	if affected != len(ids) {
		return ErrInvitationGeneric()
	}

	return nil
}
