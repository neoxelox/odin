package session

import (
	"context"
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/internal/utility"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
)

type CreatorUsecase struct {
	class.Usecase
	database          database.Database
	sessionRepository repository.SessionRepository
	userRepository    repository.UserRepository
}

func NewCreatorUsecase(configuration internal.Configuration, logger core.Logger, database database.Database,
	sessionRepository repository.SessionRepository, userRepository repository.UserRepository) *CreatorUsecase {
	return &CreatorUsecase{
		Usecase:           *class.NewUsecase(configuration, logger),
		database:          database,
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
	}
}

func (self *CreatorUsecase) Create(ctx context.Context, user model.User, metadata model.SessionMetadata) (*model.Session, *model.User, error) {
	var err error

	session := model.NewSession()
	session.UserID = user.ID
	session.Metadata = metadata

	oldSessionID := utility.CopyString(user.LastSessionID)
	user.LastSessionID = &session.ID

	err = self.database.Transaction(ctx, func(ctx context.Context) error {
		if oldSessionID != nil {
			now := time.Now()
			err = self.sessionRepository.UpdateExpiredAt(ctx, *oldSessionID, &now)
			if err != nil {
				return ErrGeneric().Wrap(err)
			}
		}

		session, err = self.sessionRepository.Create(ctx, *session)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		err = self.userRepository.UpdateSession(ctx, user.ID, *user.LastSessionID)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	return session, &user, nil
}
