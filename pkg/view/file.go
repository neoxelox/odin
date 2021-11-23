package view

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/payload"
	"github.com/neoxelox/odin/pkg/usecase/file"
	"github.com/rs/xid"
)

type FileView struct {
	class.View
	fileCreator file.CreatorUsecase
	fileGetter  file.GetterUsecase
}

func NewFileView(configuration internal.Configuration, logger core.Logger, fileCreator file.CreatorUsecase,
	fileGetter file.GetterUsecase) *FileView {
	return &FileView{
		View:        *class.NewView(configuration, logger),
		fileCreator: fileCreator,
		fileGetter:  fileGetter,
	}
}

func (self *FileView) PostFile() (interface{}, func(ctx echo.Context) error) {
	response := &payload.PostFileResponse{}
	return nil, func(ctx echo.Context) error {
		file, err := ctx.FormFile("file")
		if err != nil {
			return internal.ExcClientGeneric.Cause(err)
		}

		src, err := file.Open()
		if err != nil {
			return internal.ExcClientGeneric.Cause(err)
		}
		defer src.Close()

		fileName := fmt.Sprintf("%s%s", xid.New().String(), filepath.Ext(file.Filename))
		filePath := fmt.Sprintf("%s/%s", internal.FILES_PATH, fileName)

		dst, err := os.Create(filePath)
		if err != nil {
			return internal.ExcServerGeneric.Cause(err)
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return internal.ExcServerGeneric.Cause(err)
		}

		fileURL, err := self.fileCreator.Create(ctx.Request().Context(), fileName)
		if err != nil {
			return internal.ExcServerGeneric.Cause(err)
		}

		response.URL = fileURL

		return ctx.JSON(http.StatusOK, response)
	}
}

func (self *FileView) GetFile() (*payload.GetFileRequest, func(ctx echo.Context) error) {
	request := &payload.GetFileRequest{}
	return request, func(ctx echo.Context) error {
		filePath, err := self.fileGetter.Get(ctx.Request().Context(), request.Name)
		if err != nil {
			return internal.ExcServerGeneric.Cause(err)
		}

		return ctx.File(filePath)
	}
}
