package core

import (
	"bytes"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/labstack/echo/v4"

	"github.com/neoxelox/odin/internal"
)

var ErrRendererGeneric = internal.NewError("Renderer failed")

var TEMPLATE_EXTENSIONS = regexp.MustCompile(`^.*\.(html|txt|md)$`)

type Renderer struct {
	configuration internal.Configuration
	logger        Logger
	renderer      *template.Template
}

func NewRenderer(configuration internal.Configuration, logger Logger) *Renderer {
	logger.SetLogger(logger.Logger().With().Str("layer", "renderer").Logger())

	renderer := template.New("")
	templatesPath := filepath.Clean(internal.TEMPLATES_PATH)

	err := filepath.WalkDir(templatesPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return ErrRendererGeneric().Wrap(err)
		}

		if info.IsDir() {
			return nil
		}

		name := path[len(templatesPath)+1:]

		if !TEMPLATE_EXTENSIONS.MatchString(name) {
			return nil
		}

		file, err := ioutil.ReadFile(path)
		if err != nil {
			return ErrRendererGeneric().Wrap(err)
		}

		_, err = renderer.New(name).Parse(string(file))
		if err != nil {
			return ErrRendererGeneric().Wrap(err)
		}

		return nil
	})
	if err != nil {
		panic(ErrRendererGeneric().Wrap(err))
	}

	return &Renderer{
		configuration: configuration,
		logger:        logger,
		renderer:      renderer,
	}
}

func (self *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := self.renderer.ExecuteTemplate(w, name, data)
	if err != nil {
		return ErrRendererGeneric().Wrap(err)
	}

	return nil
}

func (self *Renderer) RenderStd(w io.Writer, template string, data interface{}) error {
	err := self.renderer.ExecuteTemplate(w, template, data)
	if err != nil {
		return ErrRendererGeneric().Wrap(err)
	}

	return nil
}

func (self *Renderer) RenderBytes(template string, data interface{}) ([]byte, error) {
	var w bytes.Buffer

	err := self.RenderStd(&w, template, data)
	if err != nil {
		return nil, ErrRendererGeneric().Wrap(err)
	}

	return w.Bytes(), nil
}

func (self *Renderer) RenderString(template string, data interface{}) (string, error) {
	bytes, err := self.RenderBytes(template, data)
	if err != nil {
		return "", ErrRendererGeneric().Wrap(err)
	}

	return string(bytes), nil
}
