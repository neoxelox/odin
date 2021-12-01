package community

import (
	"context"
	"sort"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/utility"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
)

type GetterUsecase struct {
	class.Usecase
	communityRepository  repository.CommunityRepository
	membershipRepository repository.MembershipRepository
}

func NewGetterUsecase(configuration internal.Configuration, logger core.Logger, communityRepository repository.CommunityRepository,
	membershipRepository repository.MembershipRepository) *GetterUsecase {
	return &GetterUsecase{
		Usecase:              *class.NewUsecase(configuration, logger),
		communityRepository:  communityRepository,
		membershipRepository: membershipRepository,
	}
}

func (self *GetterUsecase) Get(ctx context.Context, user model.User, communityID string) (*model.Community, *model.Membership, error) {
	community, err := self.communityRepository.GetByID(ctx, communityID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if community == nil {
		return nil, nil, ErrInvalid()
	}

	membership, err := self.membershipRepository.GetByUserAndCommunity(ctx, user.ID, community.ID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if membership == nil {
		return nil, nil, ErrNotBelongs()
	}

	if membership.DeletedAt != nil {
		return nil, nil, ErrNotBelongs()
	}

	return community, membership, nil
}

func (self *GetterUsecase) List(ctx context.Context, user model.User) ([]model.Community, []model.Membership, error) {
	memberships, err := self.membershipRepository.List(ctx, user.ID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	communityIDs := []string{}
	for _, membership := range memberships {
		communityIDs = append(communityIDs, membership.CommunityID)
	}

	communities, err := self.communityRepository.GetByIDs(ctx, communityIDs)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	sort.Slice(memberships, func(i, j int) bool {
		return memberships[i].CreatedAt.Before(memberships[j].CreatedAt)
	})

	utility.EqualSort(memberships, communities, func(i, j int) bool {
		return memberships[i].CommunityID == communities[j].ID
	})

	return communities, memberships, nil
}
