package core

import (
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"

	"github.com/neoxelox/odin/internal"
)

var ErrRendererGeneric = internal.NewError("Renderer failed")

type Renderer struct {
	configuration internal.Configuration
	logger        Logger
	renderer      *template.Template
}

func NewRenderer(configuration internal.Configuration, logger Logger) *Renderer {
	logger.SetLogger(logger.Logger().With().Str("layer", "renderer").Logger())

	glob := fmt.Sprintf("%s/*.html", internal.TEMPLATES_PATH)

	return &Renderer{
		configuration: configuration,
		logger:        logger,
		renderer:      template.Must(template.ParseGlob(glob)),
	}
}

func (self *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := self.renderer.ExecuteTemplate(w, name, data)
	if err != nil {
		return ErrRendererGeneric().Wrap(err)
	}

	return nil
}
