package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/pkg"
	"github.com/neoxelox/odin/pkg/command"
)

func main() {
	configuration := internal.NewConfiguration(nil)
	logger := core.NewLogger(*configuration)

	// cli, err := pkg.NewCLI(*configuration, *logger)
	// if err != nil {
	// 	logger.Panic(err)
	// }

	// REMOVE THIS WIP...

	database, err := database.New(context.TODO(), 5, *configuration, *logger)
	if err != nil {
		logger.Panic(err)
	}

	err = command.NewSeedCommand(*configuration, *logger, *database).Execute()
	if err != nil {
		logger.Error(err)
	}

	err = database.Close(context.TODO())
	if err != nil {
		logger.Error(err)
	}

	err = logger.Close(context.TODO())
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	// done := make(chan bool, 1)

	// go listenAndClose(&done, configuration, logger, cli)

	// // TODO: CHANGE THIS FOR NON DAEMON THINGS....

	// err = cli.Start()
	// if err != nil {
	// 	logger.Error(err)
	// 	myself, _ := os.FindProcess(syscall.Getpid())
	// 	_ = myself.Signal(os.Interrupt)
	// }

	// // Wait until all things have been gracefully closed
	// <-done
	// // TODO: EXIT FAIL IF ERROR
}

func listenAndClose(done *chan bool, configuration *internal.Configuration, logger *core.Logger, cli *pkg.CLI) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	signal := <-quit

	logger.Infof("Received signal: %v", signal)

	deadline := time.Duration(configuration.GracefulTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	err := cli.Close(ctx)
	if err != nil {
		logger.Error(err)
	}

	err = logger.Close(ctx)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	close(*done)
}
