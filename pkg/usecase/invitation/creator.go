package invitation

import (
	"context"
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
	"github.com/neoxelox/odin/pkg/usecase/community"
	"github.com/neoxelox/odin/pkg/usecase/user"
	"github.com/nyaruka/phonenumbers"
)

type CreatorUsecase struct {
	class.Usecase
	database             database.Database
	invitationRepository repository.InvitationRepository
	membershipRepository repository.MembershipRepository
	userRepository       repository.UserRepository
}

func NewCreatorUsecase(configuration internal.Configuration, logger core.Logger, database database.Database,
	invitationRepository repository.InvitationRepository, membershipRepository repository.MembershipRepository,
	userRepository repository.UserRepository) *CreatorUsecase {
	return &CreatorUsecase{
		Usecase:              *class.NewUsecase(configuration, logger),
		database:             database,
		invitationRepository: invitationRepository,
		membershipRepository: membershipRepository,
		userRepository:       userRepository,
	}
}

func (self *CreatorUsecase) Create(ctx context.Context, inviter model.User, communityID string, phone string, door string, role string) (*model.Invitation, error) {
	ph, err := phonenumbers.Parse(phone, "ES")
	if err != nil {
		return nil, user.ErrInvalidPhone().Wrap(err)
	}

	if !phonenumbers.IsValidNumber(ph) {
		return nil, user.ErrInvalidPhone()
	}

	phone = phonenumbers.Format(ph, phonenumbers.E164)

	if phone == inviter.Phone {
		return nil, ErrInvitingYourself()
	}

	if len(door) < model.MEMBERSHIP_DOOR_MIN_LENGTH || len(door) > model.MEMBERSHIP_DOOR_MAX_LENGTH {
		return nil, community.ErrInvalidDoor()
	}

	if !model.MembershipRole.Has(role) {
		return nil, community.ErrInvalidRole()
	}

	inviterMembership, err := self.membershipRepository.GetByUserAndCommunity(ctx, inviter.ID, communityID)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	if inviterMembership == nil || inviterMembership.DeletedAt != nil {
		return nil, community.ErrNotBelongs()
	}

	if inviterMembership.Role != model.MembershipRole.ADMINISTRATOR {
		return nil, community.ErrNotPermission()
	}

	existingInvitation, err := self.invitationRepository.GetByPhone(ctx, phone)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	if existingInvitation != nil && time.Now().Before(existingInvitation.ExpiresAt) {
		return nil, ErrAlreadyInvited()
	}

	invitedUser, err := self.userRepository.GetByPhone(ctx, phone)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	if invitedUser != nil {
		invitedMembership, err := self.membershipRepository.GetByUserAndCommunity(ctx, invitedUser.ID, communityID)
		if err != nil {
			return nil, ErrGeneric().Wrap(err)
		}

		if invitedMembership != nil && invitedMembership.DeletedAt == nil {
			return nil, community.ErrAlreadyJoined()
		}
	}

	invitation := model.NewInvitation()
	invitation.Phone = phone
	invitation.CommunityID = communityID
	invitation.Door = door
	invitation.Role = role

	err = self.database.Transaction(ctx, func(ctx context.Context) error {
		var err error

		if existingInvitation != nil {
			err = self.invitationRepository.DeleteByID(ctx, existingInvitation.ID)
			if err != nil {
				return ErrGeneric().Wrap(err)
			}
		}

		invitation, err = self.invitationRepository.Create(ctx, *invitation)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		return nil
	})
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	return invitation, nil
}
