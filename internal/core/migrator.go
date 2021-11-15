package core

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/neoxelox/odin/internal"
)

const POSTGRES_DSN = "postgresql://%s:%s@%s:%d/%s?sslmode=%s&x-multi-statement=true"

var ErrMigratorGeneric = internal.NewError("Migrator failed")

type Migrator struct {
	configuration internal.Configuration
	logger        Logger
	migrator      migrate.Migrate
}

func NewMigrator(configuration internal.Configuration, logger Logger) (*Migrator, error) {
	logger.SetLogger(logger.Logger().With().Str("layer", "migrator").Logger())

	path := fmt.Sprintf("file://%s", internal.MIGRATIONS_PATH)

	dsn := fmt.Sprintf(
		POSTGRES_DSN,
		configuration.DatabaseUser,
		configuration.DatabasePassword,
		configuration.DatabaseHost,
		configuration.DatabasePort,
		configuration.DatabaseName,
		configuration.DatabaseSSLMode,
	)

	migrator, err := migrate.New(path, dsn)
	if err != nil {
		return nil, ErrMigratorGeneric().Wrap(err)
	}

	migrator.Log = *NewMigrateLogger(logger)
	migrator.LockTimeout = time.Duration(configuration.GracefulTimeout) * time.Second

	return &Migrator{
		configuration: configuration,
		logger:        logger,
		migrator:      *migrator,
	}, nil
}

func (self *Migrator) Apply(ctx context.Context) error {
	if _, _, err := self.migrator.Version(); err == migrate.ErrNilVersion {
		err = self.migrator.Force(0)
		if err != nil {
			errC := self.Close(ctx)
			err = errors.CombineErrors(err, errC)
			return ErrMigratorGeneric().Wrap(err)
		}
	}

	err := self.migrator.Up()
	switch err {
	case nil:
		self.logger.Info("Applied all migrations successfully")
	case migrate.ErrNoChange:
		self.logger.Info("No migrations to apply")
	default:
		errC := self.Close(ctx)
		err := errors.CombineErrors(err, errC)
		return ErrMigratorGeneric().Wrap(err)
	}

	err = self.Close(ctx)
	if err != nil {
		return ErrMigratorGeneric().Wrap(err)
	}

	return nil
}

// TODO: opposite of apply (down)

func (self *Migrator) Close(ctx context.Context) error {
	self.logger.Info("Closing migrator")

	err, errD := self.migrator.Close()
	if err != nil {
		err = errors.CombineErrors(err, errD)
		return ErrMigratorGeneric().Wrap(err)
	}

	return nil
}

type MigrateLogger struct {
	logger Logger
}

func NewMigrateLogger(logger Logger) *MigrateLogger {
	return &MigrateLogger{
		logger: logger,
	}
}

func (self MigrateLogger) Printf(format string, v ...interface{}) {
	self.logger.Infof(format, v...) // TODO: Make logs prettier
}

func (self MigrateLogger) Verbose() bool {
	return self.logger.Verbose()
}
