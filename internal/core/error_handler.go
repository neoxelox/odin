package core

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/neoxelox/odin/internal"
)

var ErrErrorHandlerGeneric = internal.NewError("Error handler failed")

type ErrorHandler struct {
	configuration internal.Configuration
	logger        Logger
}

func NewErrorHandler(configuration internal.Configuration, logger Logger) *ErrorHandler {
	logger.SetLogger(logger.Logger().With().Str("layer", "error_handler").Logger())

	return &ErrorHandler{
		configuration: configuration,
		logger:        logger,
	}
}

func (self *ErrorHandler) Handle(err error, ctx echo.Context) {
	var exc *internal.Exception

	exc, ok := err.(*internal.Exception)
	if !ok {
		switch err {
		case echo.ErrNotFound:
			exc = internal.ExcNotFound.Cause(err)
		case echo.ErrStatusRequestEntityTooLarge:
			exc = internal.ExcInvalidRequest.Cause(err)
		case http.ErrHandlerTimeout:
			exc = internal.ExcRequestTimeout.Cause(err)
		default: // Fallback.
			exc = internal.ExcServerGeneric.Cause(err)
		}
	}

	if exc.Status >= http.StatusInternalServerError {
		self.logger.Error(exc.Origin)
	}

	if ctx.Response().Committed {
		return
	}

	if ctx.Request().Method == http.MethodHead {
		err = ctx.NoContent(exc.Status)
	} else {
		if self.configuration.Environment != internal.Environment.DEVELOPMENT {
			exc.Redact()
		}

		err = ctx.JSON(exc.Status, exc)
	}

	if err != nil {
		self.logger.Error(errors.Wrapf(err, "Cannot return exception %s", exc))
	}
}
