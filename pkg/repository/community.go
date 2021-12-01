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

const COMMUNITY_TABLE = "community"

var ErrCommunityGeneric = internal.NewError("Community query failed")

type CommunityRepository struct {
	class.Repository
}

func NewCommunityRepository(configuration internal.Configuration, logger core.Logger, database database.Database) *CommunityRepository {
	return &CommunityRepository{
		Repository: *class.NewRepository(configuration, logger, database),
	}
}

func (self *CommunityRepository) Create(ctx context.Context, community model.Community) (*model.Community, error) {
	var c model.Community

	query := fmt.Sprintf(`INSERT INTO "%s"
						  ("id", "address", "name", "categories", "pinned_ids", "created_at", "deleted_at")
						  VALUES ($1, $2, $3, $4, $5, $6, $7)
						  RETURNING *;`, COMMUNITY_TABLE)

	err := self.Database.Query(
		ctx, query, community.ID, community.Address, community.Name, community.Categories, community.PinnedIDs, community.CreatedAt, community.DeletedAt).Scan(&c)
	if err != nil {
		return nil, ErrCommunityGeneric().Wrap(err)
	}

	return &c, nil
}

func (self *CommunityRepository) GetByID(ctx context.Context, id string) (*model.Community, error) {
	var c model.Community

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "id" = $1;`, COMMUNITY_TABLE)

	err := self.Database.Query(ctx, query, id).Scan(&c)
	switch {
	case err == nil:
		return &c, nil
	case database.ErrNoRows().Is(err):
		return nil, nil
	default:
		return nil, ErrCommunityGeneric().Wrap(err)
	}
}

func (self *CommunityRepository) GetByIDs(ctx context.Context, ids []string) ([]model.Community, error) {
	var cs []model.Community

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "id" = ANY ($1);`, COMMUNITY_TABLE)

	err := self.Database.Query(ctx, query, ids).Scan(&cs)
	switch {
	case err == nil:
		return cs, nil
	case database.ErrNoRows().Is(err):
		return []model.Community{}, nil
	default:
		return nil, ErrCommunityGeneric().Wrap(err)
	}
}
