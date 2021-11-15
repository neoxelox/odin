package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	gommon "github.com/labstack/gommon/log"
	"github.com/rs/zerolog"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/core"
)

var ErrGeneric = internal.NewError("Server failed")

type Server struct {
	configuration internal.Configuration
	logger        core.Logger
	server        *echo.Echo
}

func New(configuration internal.Configuration, logger core.Logger, serializer core.Serializer,
	validator core.Validator, renderer core.Renderer, binder core.Binder, errorHandler core.ErrorHandler) *Server {
	logger.SetLogger(logger.Logger().With().Str("layer", "server").Logger())

	server := echo.New()

	echoLogger := NewEchoLogger(logger)

	server.HideBanner = true
	server.HidePort = true
	server.DisableHTTP2 = true
	server.Debug = configuration.Environment == internal.Environment.DEVELOPMENT
	server.Server.IdleTimeout = time.Duration(configuration.GracefulTimeout) * time.Second
	server.Server.MaxHeaderBytes = configuration.RequestHeaderMaxSize
	server.Server.ReadHeaderTimeout = time.Duration(configuration.GracefulTimeout) * time.Second
	server.Server.ReadTimeout = time.Duration(configuration.GracefulTimeout) * time.Second
	server.Server.WriteTimeout = time.Duration(configuration.GracefulTimeout) * time.Second

	server.Logger = echoLogger
	server.JSONSerializer = &serializer
	server.Binder = &binder
	server.Renderer = &renderer
	server.Validator = &validator
	server.HTTPErrorHandler = errorHandler.Handle
	server.IPExtractor = echo.ExtractIPFromRealIPHeader()

	return &Server{
		configuration: configuration,
		logger:        logger,
		server:        server,
	}
}

func (self *Server) Run() error {
	self.logger.Info("Starting server")

	err := self.server.Start(fmt.Sprintf(":%d", self.configuration.AppPort))
	if err != nil && err != http.ErrServerClosed {
		return ErrGeneric().Wrap(err)
	}

	return nil
}

func (self *Server) Use(middleware ...echo.MiddlewareFunc) {
	self.server.Pre(middleware...)
}

func (self *Server) Host(name string, middleware ...echo.MiddlewareFunc) *echo.Group {
	return self.server.Host(name, middleware...)
}

func (self *Server) Close(ctx context.Context) error {
	self.logger.Info("Closing server")

	self.server.Server.SetKeepAlivesEnabled(false)

	err := self.server.Shutdown(ctx)
	if err != nil {
		return ErrGeneric().Wrap(err)
	}

	return nil
}

type EchoLogger struct {
	core.Logger
}

func NewEchoLogger(logger core.Logger) *EchoLogger {
	return &EchoLogger{
		Logger: logger,
	}
}

func (self EchoLogger) Level() gommon.Lvl {
	return self.GLevel()
}

func (self *EchoLogger) SetLevel(v gommon.Lvl) {
	self.SetGLevel(v)
}

func (self EchoLogger) Printj(j gommon.JSON) {
	raw, _ := json.Marshal(j)
	self.Logger.Logger().Log().RawJSON(zerolog.MessageFieldName, raw).Msg("")
}

func (self EchoLogger) Debugj(j gommon.JSON) {
	raw, _ := json.Marshal(j)
	self.Logger.Logger().Debug().RawJSON(zerolog.MessageFieldName, raw).Msg("")
}

func (self EchoLogger) Infoj(j gommon.JSON) {
	raw, _ := json.Marshal(j)
	self.Logger.Logger().Info().RawJSON(zerolog.MessageFieldName, raw).Msg("")
}

func (self EchoLogger) Warnj(j gommon.JSON) {
	raw, _ := json.Marshal(j)
	self.Logger.Logger().Warn().RawJSON(zerolog.MessageFieldName, raw).Msg("")
}

func (self EchoLogger) Errorj(j gommon.JSON) {
	raw, _ := json.Marshal(j)
	self.Logger.Logger().Error().RawJSON(zerolog.MessageFieldName, raw).Msg("")
}

func (self EchoLogger) Fatalj(j gommon.JSON) {
	raw, _ := json.Marshal(j)
	self.Logger.Logger().Fatal().RawJSON(zerolog.MessageFieldName, raw).Msg("")
}

func (self EchoLogger) Panicj(j gommon.JSON) {
	raw, _ := json.Marshal(j)
	self.Logger.Logger().Panic().RawJSON(zerolog.MessageFieldName, raw).Msg("")
}
