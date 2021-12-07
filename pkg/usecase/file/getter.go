package file

import (
	"context"
	"fmt"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
)

type GetterUsecase struct {
	class.Usecase
}

func NewGetterUsecase(configuration internal.Configuration, logger core.Logger) *GetterUsecase {
	return &GetterUsecase{
		Usecase: *class.NewUsecase(configuration, logger),
	}
}

func (self *GetterUsecase) Get(ctx context.Context, fileName string) (string, error) {
	filePath := fmt.Sprintf("%s/%s", internal.FILES_PATH, fileName)

	return filePath, nil
}
