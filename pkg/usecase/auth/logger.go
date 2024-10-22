package auth

import (
	"context"
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
	"github.com/neoxelox/odin/pkg/usecase/invitation"
	"github.com/neoxelox/odin/pkg/usecase/otp"
	"github.com/neoxelox/odin/pkg/usecase/session"
	"github.com/neoxelox/odin/pkg/usecase/user"
)

type LoggerUsecase struct {
	class.Usecase
	database          database.Database
	otpVerifier       otp.VerifierUsecase
	userCreator       user.CreatorUsecase
	sessionCreator    session.CreatorUsecase
	authCreator       CreatorUsecase
	invitationCreator invitation.CreatorUsecase // TODO: REMOVE TEMPORAL INVITATION
	otpRepository     repository.OTPRepository
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
}

func NewLoggerUsecase(configuration internal.Configuration, logger core.Logger, database database.Database,
	otpVerifier otp.VerifierUsecase, userCreator user.CreatorUsecase, sessionCreator session.CreatorUsecase,
	authCreator CreatorUsecase, otpRepository repository.OTPRepository, userRepository repository.UserRepository,
	sessionRepository repository.SessionRepository, invitationCreator invitation.CreatorUsecase) *LoggerUsecase {
	return &LoggerUsecase{
		Usecase:           *class.NewUsecase(configuration, logger),
		database:          database,
		otpVerifier:       otpVerifier,
		userCreator:       userCreator,
		sessionCreator:    sessionCreator,
		authCreator:       authCreator,
		invitationCreator: invitationCreator,
		otpRepository:     otpRepository,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
	}
}

func (self *LoggerUsecase) Login(ctx context.Context, otpID string, code string, metadata model.SessionMetadata) (string, *model.User, error) {
	otpReq, err := self.otpVerifier.Verify(ctx, otpID, code, model.OTPType.SMS)
	if err != nil {
		if !otp.ErrGeneric().Is(err) {
			return "", nil, ErrGeneric().As(err)
		}

		return "", nil, ErrGeneric().Wrap(err)
	}

	var accessToken string
	var user *model.User
	err = self.database.Transaction(ctx, func(ctx context.Context) error {
		err := self.otpRepository.Delete(ctx, otpReq.ID)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		user, err = self.userRepository.GetByPhone(ctx, otpReq.Asset)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		if user == nil {
			user, err = self.userCreator.Create(ctx, otpReq.Asset)
			if err != nil {
				return ErrGeneric().Wrap(err)
			}

			// TODO: REMOVE TEMPORAL INVITATION
			admin, err := self.userRepository.GetByID(ctx, "9bsv0s7q8b4002uqbcng")
			if err != nil {
				return ErrGeneric().Wrap(err)
			}
			_, err = self.invitationCreator.Create(ctx, *admin, "9bsv0s5a5rsg02purd40", user.Phone, "5º 1ª", model.MembershipRole.RESIDENT)
			if err != nil {
				return ErrGeneric().Wrap(err)
			}
		}

		var session *model.Session

		session, user, err = self.sessionCreator.Create(ctx, *user, metadata)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		accessToken, err = self.authCreator.Create(ctx, *session)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		return nil
	})
	if err != nil {
		return "", nil, ErrGeneric().Wrap(err)
	}

	if user.Name == "" {
		return accessToken, nil, nil
	}

	return accessToken, user, nil
}

func (self *LoggerUsecase) Logout(ctx context.Context, session model.Session) error {
	err := self.sessionRepository.UpdateExpiredAt(ctx, session.ID, time.Now())
	if err != nil {
		return ErrGeneric().Wrap(err)
	}

	return nil
}
