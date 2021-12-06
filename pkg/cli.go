package pkg

import (
	"context"

	"github.com/mkideal/cli"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/cache"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/pkg/command"
)

type CLI struct {
	class.CLI
}

func NewCLI(configuration internal.Configuration, logger core.Logger) (*CLI, error) {
	ctx := context.Background()

	/* DEPENDENCIES */

	database, err := database.New(ctx, 5, configuration, logger)
	if err != nil {
		return nil, ErrAPIGeneric().Wrap(err)
	}

	cache, err := cache.New(ctx, 5, configuration, logger)
	if err != nil {
		return nil, ErrAPIGeneric().Wrap(err)
	}

	/* MIDDLEWARES */

	/* SERVICES */

	/* REPOSITORIES  */

	/* USECASES */

	/* COMMANDS */

	helpCommand := cli.HelpCommand("display help information")
	_ = command.NewSeedCommand(configuration, logger, *database)

	root := cli.Root(
		helpCommand,
		cli.Tree(helpCommand),
		//cli.Tree(seedCommand.Execute()), // TODO: THIS IS NOSENSE, CHANGE ALL OF THIS!
	)

	return &CLI{
		CLI: *class.NewCLI(
			// ON START
			func() error {
				return nil
			},
			// ON CLOSE
			func(ctx context.Context) error {
				return nil
			}, configuration, logger, *database, *cache, root),
	}, nil
}
