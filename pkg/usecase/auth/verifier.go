package auth

import (
	"context"
	"strings"
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/utility"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/repository"
	"github.com/vk-rv/pvx"
)

type VerifierUsecase struct {
	class.Usecase
	key               *pvx.SymKey
	codifier          *pvx.ProtoV4Local
	sessionRepository repository.SessionRepository
	userRepository    repository.UserRepository
}

func NewVerifierUsecase(configuration internal.Configuration, logger core.Logger, sessionRepository repository.SessionRepository,
	userRepository repository.UserRepository) *VerifierUsecase {
	return &VerifierUsecase{
		Usecase:           *class.NewUsecase(configuration, logger),
		key:               pvx.NewSymmetricKey([]byte(configuration.SessionKey), pvx.Version4),
		codifier:          pvx.NewPV4Local(),
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
	}
}

func (self *VerifierUsecase) Verify(ctx context.Context, accessToken string, path string) (*model.Session, *model.User, error) {
	decoded := &model.AccessToken{}

	err := self.codifier.Decrypt(accessToken, self.key).Scan(&decoded.Private, &decoded.Public)
	if err != nil {
		return nil, nil, ErrTamperedAccessToken().Wrap(err)
	}

	if time.Now().After(decoded.Private.ExpiresAt) {
		return nil, nil, ErrExpiredAccessToken()
	}

	if !utility.StringIn(path, UNVERSIONED_PATHS) && strings.Split(path, "/")[1] != decoded.Public.ApiVersion {
		return nil, nil, ErrInvalidAccessToken()
	}

	session, err := self.sessionRepository.GetByID(ctx, decoded.Private.SessionID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if session == nil {
		return nil, nil, ErrInvalidAccessToken()
	}

	if session.ExpiredAt != nil && time.Now().After(*session.ExpiredAt) {
		return nil, nil, ErrExpiredSession()
	}

	user, err := self.userRepository.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	if user == nil {
		return nil, nil, ErrInvalidAccessToken()
	}

	if *user.LastSessionID != session.ID {
		return nil, nil, ErrInvalidAccessToken()
	}

	if user.IsBanned {
		return nil, nil, ErrBannedUser()
	}

	session.LastSeenAt = time.Now()

	err = self.sessionRepository.UpdateLastSeen(ctx, session.ID, session.LastSeenAt)
	if err != nil {
		return nil, nil, ErrGeneric().Wrap(err)
	}

	return session, user, nil
}
