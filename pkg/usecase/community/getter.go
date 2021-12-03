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
	userRepository       repository.UserRepository
}

func NewGetterUsecase(configuration internal.Configuration, logger core.Logger, communityRepository repository.CommunityRepository,
	membershipRepository repository.MembershipRepository, userRepository repository.UserRepository) *GetterUsecase {
	return &GetterUsecase{
		Usecase:              *class.NewUsecase(configuration, logger),
		communityRepository:  communityRepository,
		membershipRepository: membershipRepository,
		userRepository:       userRepository,
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
	memberships, err := self.membershipRepository.ListByUser(ctx, user.ID)
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

func (self *GetterUsecase) GetUser(ctx context.Context, requester model.User, communityID string, membershipID string) (*model.User, *model.Membership, error) {
	requesterMembership, err := self.membershipRepository.GetByUserAndCommunity(ctx, requester.ID, communityID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if requesterMembership == nil || requesterMembership.DeletedAt != nil {
		return nil, nil, ErrNotBelongs()
	}

	if requesterMembership.ID == membershipID {
		return &requester, requesterMembership, nil
	}

	membership, err := self.membershipRepository.GetByID(ctx, membershipID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if membership == nil {
		return nil, nil, ErrNotBelongs()
	}

	if membership.CommunityID != requesterMembership.CommunityID {
		return nil, nil, ErrNotBelongs()
	}

	user, err := self.userRepository.GetByID(ctx, membership.UserID)
	if err != nil {
		return nil, nil, ErrGeneric()
	}

	if user == nil {
		return nil, nil, ErrInvalid()
	}

	return user, membership, nil
}

func (self *GetterUsecase) ListUsers(ctx context.Context, user model.User, communityID string) ([]model.User, []model.Membership, error) {
	membership, err := self.membershipRepository.GetByUserAndCommunity(ctx, user.ID, communityID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if membership == nil || membership.DeletedAt != nil {
		return nil, nil, ErrNotBelongs()
	}

	memberships, err := self.membershipRepository.ListByCommunity(ctx, communityID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	userIDs := []string{}
	for _, membership := range memberships {
		userIDs = append(userIDs, membership.UserID)
	}

	users, err := self.userRepository.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Name < users[j].Name
	})

	utility.EqualSort(users, memberships, func(i, j int) bool {
		return users[i].ID == memberships[j].UserID
	})

	return users, memberships, nil
}
