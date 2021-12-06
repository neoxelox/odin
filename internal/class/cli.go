package class

import (
	"context"
	"os"

	"github.com/mkideal/cli"

	"github.com/cockroachdb/errors"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/cache"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
)

var ErrCLIGeneric = internal.NewError("CLI failed")

type CLI struct {
	start         func() error
	close         func(context.Context) error
	Configuration internal.Configuration
	Logger        core.Logger
	Database      database.Database
	Cache         cache.Cache
	Root          *cli.Command // CHANGE THIS, CREATE AN INTERNAL RUNNER
}

func NewCLI(start func() error, close func(context.Context) error, configuration internal.Configuration,
	logger core.Logger, database database.Database, cache cache.Cache, root *cli.Command) *CLI {
	logger.SetLogger(logger.Logger().With().Str("layer", "cli").Logger())

	return &CLI{
		start:         start,
		close:         close,
		Configuration: configuration,
		Logger:        logger,
		Database:      database,
		Cache:         cache,
		Root:          root,
	}
}

func (self *CLI) Start() error {
	self.Logger.Info("Starting cli")

	err := self.start()
	if err != nil {
		return ErrCLIGeneric().Wrap(err)
	}

	err = self.Root.Run(os.Args[1:])
	if err != nil {
		return ErrCLIGeneric().Wrap(err)
	}

	return nil
}

func (self *CLI) Close(ctx context.Context) error {
	self.Logger.Info("Closing cli")

	var err error

	// TODO CLOSE ROOT

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
		return ErrCLIGeneric().Wrap(err)
	}

	return nil
}
