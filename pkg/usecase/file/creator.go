package file

import (
	"context"
	"fmt"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
)

type CreatorUsecase struct {
	class.Usecase
}

func NewCreatorUsecase(configuration internal.Configuration, logger core.Logger) *CreatorUsecase {
	return &CreatorUsecase{
		Usecase: *class.NewUsecase(configuration, logger),
	}
}

func (self *CreatorUsecase) Create(ctx context.Context, fileName string) (string, error) {
	fileURL := fmt.Sprintf("http://%s:%d/file/%s", self.Configuration.AppHost, self.Configuration.AppPort, fileName)

	if self.Configuration.Environment == internal.Environment.PRODUCTION {
		return "", ErrGeneric()
	}

	return fileURL, nil
}
