package class

import (
	"github.com/mkideal/cli"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/core"
)

type Command struct {
	Configuration internal.Configuration
	Logger        core.Logger
}

// TODO: CHANGE ALL OF THIS AND RETURN A CLI COMMAND WTF!!!
func NewCommand(configuration internal.Configuration, logger core.Logger) *Command {
	logger.SetLogger(logger.Logger().With().Str("layer", "command").Logger())

	return &Command{
		Configuration: configuration,
		Logger:        logger,
	}
}

type CommandHandler func(ctx *cli.Context) error

type CommandEndpoint struct {
	Name        string
	Aliases     []string
	Description string
	Details     string
	Arguments   interface{}
}

func (self *Command) Handle(endpoint CommandEndpoint, handler CommandHandler) *cli.Command {
	return &cli.Command{
		Name:    endpoint.Name,
		Aliases: endpoint.Aliases,
		Desc:    endpoint.Description,
		Text:    endpoint.Details,
		Argv:    func() interface{} { return endpoint.Arguments },
		Fn: func(ctx *cli.Context) error {
			// TODO: BINDING
			// TODO: VALIDATE
			// TODO: ERROR HANDLING
			return handler(ctx)
		},
	}
}
