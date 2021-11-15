package core

import (
	"github.com/labstack/echo/v4"

	"github.com/neoxelox/odin/internal"
)

var ErrSerializerGeneric = internal.NewError("Serializer failed")

type Serializer struct {
	configuration internal.Configuration
	logger        Logger
	serializer    echo.DefaultJSONSerializer
}

func NewSerializer(configuration internal.Configuration, logger Logger) *Serializer {
	logger.SetLogger(logger.Logger().With().Str("layer", "serializer").Logger())

	return &Serializer{
		configuration: configuration,
		logger:        logger,
		serializer:    echo.DefaultJSONSerializer{},
	}
}

func (self *Serializer) Serialize(c echo.Context, i interface{}, indent string) error {
	err := self.serializer.Serialize(c, i, indent)
	if err != nil {
		return ErrSerializerGeneric().Wrap(err)
	}

	return nil
}

func (self *Serializer) Deserialize(c echo.Context, i interface{}) error {
	err := self.serializer.Deserialize(c, i)
	if err != nil {
		return ErrSerializerGeneric().Wrap(err)
	}

	return nil
}
