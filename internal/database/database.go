package database

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/core"
	"github.com/randallmlough/pgxscan"
)

const (
	POSTGRES_DSN                                = "postgresql://%s:%s@%s:%d/%s?sslmode=%s"
	CONTEXT_TRANSACTION_KEY internal.ContextKey = "database:tx"
)

var (
	ErrGeneric     = internal.NewError("Database failed")
	ErrBelowMin    = internal.NewError("Database pool size below minimum")
	ErrTransaction = internal.NewError("Database transaction failed")
)

type Database struct {
	configuration internal.Configuration
	logger        core.Logger
	pool          *pgxpool.Pool
}

func New(ctx context.Context, retries int, configuration internal.Configuration, logger core.Logger) (*Database, error) {
	logger.SetLogger(logger.Logger().With().Str("layer", "database").Logger())

	dsn := fmt.Sprintf(
		POSTGRES_DSN,
		configuration.DatabaseUser,
		configuration.DatabasePassword,
		configuration.DatabaseHost,
		configuration.DatabasePort,
		configuration.DatabaseName,
		configuration.DatabaseSSLMode,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	config.MinConns = int32(configuration.DatabaseMinConns)
	config.MaxConns = int32(configuration.DatabaseMaxConns)
	config.ConnConfig.RuntimeParams["standard_conforming_strings"] = "on"
	config.ConnConfig.RuntimeParams["application_name"] = configuration.AppName

	pgxLogger := NewPgxLogger(logger)
	pgxLogLevel := pgxLogger.logger.PLevel()

	// PGX Info level is too much!
	if pgxLogLevel == pgx.LogLevelInfo {
		pgxLogLevel = pgx.LogLevelError
	}

	config.ConnConfig.Logger = pgxLogger
	config.ConnConfig.LogLevel = pgxLogLevel

	migrator, err := core.NewMigrator(configuration, logger)
	if err != nil {
		return nil, ErrGeneric().Wrap(err)
	}

	delay := time.NewTicker(1 * time.Second)
	timeoutExceeded := time.After((time.Duration(retries) * time.Second))

	// TODO: Instead of while do --> do while in order not to waste 1 second
	for {
		select {
		case <-timeoutExceeded:
			return nil, ErrGeneric()
		case <-delay.C:
			logger.Info("Trying to connect to the database")

			pool, err := pgxpool.ConnectConfig(ctx, config)
			if err == nil {
				logger.Info("Connected to the database")

				err = migrator.Apply(ctx)
				if err != nil {
					pool.Close()
					return nil, ErrGeneric().Wrap(err)
				}

				return &Database{
					configuration: configuration,
					logger:        logger,
					pool:          pool,
				}, nil
			}
		}
	}
}

func (self *Database) Health(ctx context.Context) error {
	delay := time.NewTicker(100 * time.Millisecond)
	timeoutExceeded := time.After(300 * time.Millisecond)

	for {
		select {
		case <-timeoutExceeded:
			return ErrGeneric()
		case <-delay.C:
			err := func() error {
				if self.pool.Stat().TotalConns() < int32(self.configuration.DatabaseMinConns) {
					return ErrBelowMin()
				}

				if err := self.pool.Ping(ctx); err != nil {
					return err
				}

				if err := ctx.Err(); err != nil {
					return err
				}

				rows, err := self.pool.Query(ctx, `SELECT true;`)
				if err != nil {
					return err
				}

				var ok bool
				if err = pgxscan.NewScanner(rows).Scan(&ok); !ok || err != nil {
					return err
				}

				if err := ctx.Err(); err != nil {
					return err
				}

				return nil
			}()

			if err != nil {
				return ErrGeneric().Wrap(err)
			}

			return nil
		}
	}
}

type Scan struct {
	Scan func(dest ...interface{}) error
}

func (self *Database) Query(ctx context.Context, sql string, args ...interface{}) Scan {
	return Scan{
		Scan: func(dest ...interface{}) error {
			var rows pgx.Rows
			var err error

			if ctx.Value(CONTEXT_TRANSACTION_KEY) != nil {
				transaction, _ := ctx.Value(CONTEXT_TRANSACTION_KEY).(pgx.Tx)
				rows, err = transaction.Query(ctx, sql, args...)
			} else {
				rows, err = self.pool.Query(ctx, sql, args...)
			}

			if err != nil {
				return internalError(err)
			}

			if err := ctx.Err(); err != nil {
				return internalError(err)
			}

			if err := pgxscan.NewScanner(rows).Scan(dest...); err != nil {
				return internalError(err)
			}

			return nil
		},
	}
}

func (self *Database) Exec(ctx context.Context, sql string, args ...interface{}) (int, error) {
	var command pgconn.CommandTag
	var err error

	if ctx.Value(CONTEXT_TRANSACTION_KEY) != nil {
		transaction, _ := ctx.Value(CONTEXT_TRANSACTION_KEY).(pgx.Tx)
		command, err = transaction.Exec(ctx, sql, args...)
	} else {
		command, err = self.pool.Exec(ctx, sql, args...)
	}

	if err != nil {
		return 0, internalError(err)
	}

	if err := ctx.Err(); err != nil {
		return 0, internalError(err)
	}

	return int(command.RowsAffected()), nil
}

func (self *Database) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if ctx.Value(CONTEXT_TRANSACTION_KEY) != nil {
		err := fn(ctx)
		if err != nil {
			return ErrTransaction().WrapWithDepth(1, err)
		}
		return nil
	}

	transaction, err := self.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return ErrTransaction().WrapWithDepth(1, err)
	}

	if err := ctx.Err(); err != nil {
		return ErrTransaction().WrapWithDepth(1, err)
	}

	defer func() {
		if pan := recover(); pan != nil {
			_ = transaction.Rollback(ctx)
			panic(pan)
		}
	}()

	err = fn(context.WithValue(ctx, CONTEXT_TRANSACTION_KEY, transaction))
	if err != nil {
		errR := transaction.Rollback(ctx)
		err := errors.CombineErrors(err, errR)
		return ErrTransaction().WrapWithDepth(1, err)
	}

	if err := ctx.Err(); err != nil {
		errR := transaction.Rollback(ctx)
		err := errors.CombineErrors(err, errR)
		return ErrTransaction().WrapWithDepth(1, err)
	}

	if err := transaction.Commit(ctx); err != nil {
		errR := transaction.Rollback(ctx)
		err := errors.CombineErrors(err, errR)
		return ErrTransaction().WrapWithDepth(1, err)
	}

	if err := ctx.Err(); err != nil {
		errR := transaction.Rollback(ctx)
		err := errors.CombineErrors(err, errR)
		return ErrTransaction().WrapWithDepth(1, err)
	}

	return nil
}

func (self *Database) Close(ctx context.Context) error {
	self.logger.Info("Closing database")
	self.pool.Close()

	return nil
}

type PgxLogger struct {
	logger core.Logger
}

func NewPgxLogger(logger core.Logger) *PgxLogger {
	return &PgxLogger{
		logger: logger,
	}
}

func (self PgxLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	self.logger.Logger().WithLevel(core.PlevelToZlevel[level]).Fields(data).Msg(msg)
}
