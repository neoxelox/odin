package core

import (
	"github.com/go-playground/validator/v10"
	"github.com/neoxelox/odin/internal"
)

var ErrValidatorGeneric = internal.NewError("Validator failed")

type Validator struct {
	configuration internal.Configuration
	logger        Logger
	validator     validator.Validate
}

func NewValidator(configuration internal.Configuration, logger Logger) *Validator {
	logger.SetLogger(logger.Logger().With().Str("layer", "validator").Logger())

	return &Validator{
		configuration: configuration,
		logger:        logger,
		validator:     *validator.New(),
	}
}

func (self *Validator) Validate(i interface{}) error {
	if err := self.validator.Struct(i); err != nil {
		return ErrValidatorGeneric().Wrap(err)
	}

	return nil
}
