package class

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/cache"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/internal/server"
)

var ErrAPIGeneric = internal.NewError("API failed")

type API struct {
	start         func() error
	close         func(context.Context) error
	Configuration internal.Configuration
	Logger        core.Logger
	Database      database.Database
	Cache         cache.Cache
	Server        server.Server
}

func NewAPI(start func() error, close func(context.Context) error, configuration internal.Configuration,
	logger core.Logger, database database.Database, cache cache.Cache, server server.Server) *API {
	logger.SetLogger(logger.Logger().With().Str("layer", "api").Logger())

	return &API{
		start:         start,
		close:         close,
		Configuration: configuration,
		Logger:        logger,
		Database:      database,
		Cache:         cache,
		Server:        server,
	}
}

func (self *API) Start() error {
	self.Logger.Info("Starting api")

	err := self.start()
	if err != nil {
		return ErrAPIGeneric().Wrap(err)
	}

	err = self.Server.Run()
	if err != nil {
		return ErrAPIGeneric().Wrap(err)
	}

	return nil
}

func (self *API) Close(ctx context.Context) error {
	self.Logger.Info("Closing api")

	var err error

	if errC := self.Server.Close(ctx); errC != nil {
		err = errors.CombineErrors(err, errC)
	}

	if errC := self.close(ctx); errC != nil {
		err = errors.CombineErrors(err, errC)
	}

	if errC := self.Cache.Close(ctx); errC != nil {
		err = errors.CombineErrors(err, errC)
	}

	if errC := self.Database.Close(ctx); errC != nil {
		err = errors.CombineErrors(err, errC)
	}

	if err != nil {
		return ErrAPIGeneric().Wrap(err)
	}

	return nil
}
