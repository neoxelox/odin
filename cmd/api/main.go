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
	"github.com/neoxelox/odin/pkg"
)

func main() {
	configuration := internal.NewConfiguration(nil)
	logger := core.NewLogger(*configuration)

	api, err := pkg.NewAPI(*configuration, *logger)
	if err != nil {
		logger.Panic(err)
	}

	done := make(chan bool, 1)

	go listenAndClose(&done, configuration, logger, api)

	err = api.Start()
	if err != nil {
		logger.Error(err)
		myself, _ := os.FindProcess(syscall.Getpid())
		_ = myself.Signal(os.Interrupt)
	}

	// Wait until all things have been gracefully closed
	<-done
}

func listenAndClose(done *chan bool, configuration *internal.Configuration, logger *core.Logger, api *pkg.API) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	signal := <-quit

	logger.Infof("Received signal: %v", signal)

	deadline := time.Duration(configuration.GracefulTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	err := api.Close(ctx)
	if err != nil {
		logger.Error(err)
	}

	err = logger.Close(ctx)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	close(*done)
}
