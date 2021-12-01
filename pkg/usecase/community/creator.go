package community

import (
	"context"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
	"github.com/scylladb/go-set/strset"
)

type CreatorUsecase struct {
	class.Usecase
	database            database.Database
	communityJoiner     JoinerUsecase
	communityRepository repository.CommunityRepository
}

func NewCreatorUsecase(configuration internal.Configuration, logger core.Logger, database database.Database,
	communityJoiner JoinerUsecase, communityRepository repository.CommunityRepository) *CreatorUsecase {
	return &CreatorUsecase{
		Usecase:             *class.NewUsecase(configuration, logger),
		database:            database,
		communityJoiner:     communityJoiner,
		communityRepository: communityRepository,
	}
}

func (self *CreatorUsecase) Create(ctx context.Context, address string, name *string, categories *[]string, creator *model.User) (*model.Community, *model.Membership, error) {
	community := model.NewCommunity()
	communityCategories := strset.New()

	community.Address = address

	community.Name = community.Address
	if name != nil {
		community.Name = *name
	}

	communityCategories.Add(model.COMMUNITY_DEFAULT_CATEGORIES...)
	if categories != nil {
		communityCategories.Add(*categories...)
	}
	community.Categories = communityCategories.List()

	if len(community.Address) < model.COMMUNITY_ADDRESS_MIN_LENGTH || len(community.Address) > model.COMMUNITY_ADDRESS_MAX_LENGTH {
		return nil, nil, ErrInvalidAddress()
	}

	if len(community.Name) < model.COMMUNITY_NAME_MIN_LENGTH || len(community.Name) > model.COMMUNITY_NAME_MAX_LENGTH {
		return nil, nil, ErrInvalidName()
	}

	var creatorMembership *model.Membership
	err := self.database.Transaction(ctx, func(ctx context.Context) error {
		var err error

		community, err = self.communityRepository.Create(ctx, *community)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		if creator == nil {
			return nil
		}

		creatorMembership, err = self.communityJoiner.Join(ctx, *creator, community.ID, "", model.MembershipRole.ADMINISTRATOR)
		if err != nil {
			return ErrGeneric().As(err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, ErrGeneric().As(err)
	}

	if creator == nil {
		return community, nil, nil
	}

	return community, creatorMembership, nil
}
