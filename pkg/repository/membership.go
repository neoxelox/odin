package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/pkg/model"
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

func (self *MembershipRepository) Create(ctx context.Context, membership model.Membership) (*model.Membership, error) {
	var m model.Membership

	query := fmt.Sprintf(`INSERT INTO "%s"
						  ("id", "user_id", "community_id", "door", "role", "created_at", "deleted_at")
						  VALUES ($1, $2, $3, $4, $5, $6, $7)
						  RETURNING *;`, MEMBERSHIP_TABLE)

	err := self.Database.Query(
		ctx, query, membership.ID, membership.UserID, membership.CommunityID, membership.Door, membership.Role, membership.CreatedAt, membership.DeletedAt).Scan(&m)
	if err != nil {
		return nil, ErrMembershipGeneric().Wrap(err)
	}

	return &m, nil
}

func (self *MembershipRepository) GetByID(ctx context.Context, id string) (*model.Membership, error) {
	var m model.Membership

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "id" = $1;`, MEMBERSHIP_TABLE)

	err := self.Database.Query(ctx, query, id).Scan(&m)
	switch {
	case err == nil:
		return &m, nil
	case database.ErrNoRows().Is(err):
		return nil, nil
	default:
		return nil, ErrMembershipGeneric().Wrap(err)
	}
}

func (self *MembershipRepository) GetByUserAndCommunity(ctx context.Context, userID string, communityID string) (*model.Membership, error) {
	var m model.Membership

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "user_id" = $1 AND "community_id" = $2;`, MEMBERSHIP_TABLE)

	err := self.Database.Query(ctx, query, userID, communityID).Scan(&m)
	switch {
	case err == nil:
		return &m, nil
	case database.ErrNoRows().Is(err):
		return nil, nil
	default:
		return nil, ErrMembershipGeneric().Wrap(err)
	}
}

func (self *MembershipRepository) GetByIDsAndCommunity(ctx context.Context, ids []string, communityID string) ([]model.Membership, error) {
	var ms []model.Membership

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "id" = ANY ($1) AND "community_id" = $2;`, MEMBERSHIP_TABLE)

	err := self.Database.Query(ctx, query, ids, communityID).Scan(&ms)
	switch {
	case err == nil:
		return ms, nil
	case database.ErrNoRows().Is(err):
		return []model.Membership{}, nil
	default:
		return nil, ErrMembershipGeneric().Wrap(err)
	}
}

func (self *MembershipRepository) ListByUser(ctx context.Context, userID string) ([]model.Membership, error) {
	var ms []model.Membership

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "user_id" = $1 AND "deleted_at" IS NULL;`, MEMBERSHIP_TABLE)

	err := self.Database.Query(ctx, query, userID).Scan(&ms)
	switch {
	case err == nil:
		return ms, nil
	case database.ErrNoRows().Is(err):
		return []model.Membership{}, nil
	default:
		return nil, ErrMembershipGeneric().Wrap(err)
	}
}

func (self *MembershipRepository) ListByCommunity(ctx context.Context, communityID string) ([]model.Membership, error) {
	var ms []model.Membership

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "community_id" = $1 AND "deleted_at" IS NULL;`, MEMBERSHIP_TABLE)

	err := self.Database.Query(ctx, query, communityID).Scan(&ms)
	switch {
	case err == nil:
		return ms, nil
	case database.ErrNoRows().Is(err):
		return []model.Membership{}, nil
	default:
		return nil, ErrMembershipGeneric().Wrap(err)
	}
}

func (self *MembershipRepository) UpdateDeletedAt(ctx context.Context, id string, deletedAt *time.Time) error {
	query := fmt.Sprintf(`UPDATE "%s"
						  SET "deleted_at" = $1
						  WHERE "id" = $2;`, MEMBERSHIP_TABLE)

	affected, err := self.Database.Exec(ctx, query, deletedAt, id)
	if err != nil {
		return ErrMembershipGeneric().Wrap(err)
	}

	if affected != 1 {
		return ErrMembershipGeneric()
	}

	return nil
}
