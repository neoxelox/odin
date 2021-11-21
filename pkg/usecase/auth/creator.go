package auth

import (
	"context"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/vk-rv/pvx"
)

type CreatorUsecase struct {
	class.Usecase
	key      *pvx.SymKey
	codifier *pvx.ProtoV4Local
}

func NewCreatorUsecase(configuration internal.Configuration, logger core.Logger) *CreatorUsecase {
	return &CreatorUsecase{
		Usecase:  *class.NewUsecase(configuration, logger),
		key:      pvx.NewSymmetricKey([]byte(configuration.SessionKey), pvx.Version4),
		codifier: pvx.NewPV4Local(),
	}
}

func (self *CreatorUsecase) Create(ctx context.Context, session model.Session) (string, error) {
	accessToken := model.NewAccessToken()
	accessToken.Private.SessionID = session.ID
	accessToken.Public.ApiVersion = session.Metadata.ApiVersion

	encoded, err := self.codifier.Encrypt(self.key, &accessToken.Private, pvx.WithFooter(&accessToken.Public))
	if err != nil {
		return "", ErrGeneric().Wrap(err)
	}

	return encoded, nil
}
