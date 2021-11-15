package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/core"
)

const (
	REDIS_DSN = "%s:%d"
)

var (
	ErrGeneric  = internal.NewError("Cache failed")
	ErrBelowMin = internal.NewError("Cache pool size below minimum")
)

type Cache struct {
	configuration internal.Configuration
	logger        core.Logger
	pool          *redis.Client
	cache         *cache.Cache
}

func New(ctx context.Context, retries int, configuration internal.Configuration, logger core.Logger) (*Cache, error) {
	logger.SetLogger(logger.Logger().With().Str("layer", "cache").Logger())

	redis.SetLogger(NewRedisLogger(logger))

	dsn := fmt.Sprintf(REDIS_DSN, configuration.CacheHost, configuration.CachePort)

	config := &redis.Options{
		Addr:         dsn,
		Password:     configuration.CachePassword,
		MaxRetries:   3,
		DialTimeout:  time.Duration(configuration.GracefulTimeout) * time.Second,
		ReadTimeout:  time.Duration(configuration.GracefulTimeout) * time.Second,
		PoolTimeout:  time.Duration(configuration.GracefulTimeout) * time.Second,
		WriteTimeout: time.Duration(configuration.GracefulTimeout) * time.Second,
		MinIdleConns: configuration.CacheMinConns,
		PoolSize:     configuration.CacheMaxConns,
	}

	delay := time.NewTicker(1 * time.Second)
	timeoutExceeded := time.After((time.Duration(retries) * time.Second))

	// TODO: Instead of while do --> do while in order not to waste 1 second
	for {
		select {
		case <-timeoutExceeded:
			return nil, ErrGeneric()
		case <-delay.C:
			logger.Info("Trying to connect to the cache")

			pool := redis.NewClient(config)

			result, err := pool.Ping(ctx).Result()
			if err == nil && result == "PONG" {
				logger.Info("Connected to the cache")

				cache := cache.New(&cache.Options{
					Redis:        pool,
					LocalCache:   cache.NewTinyLFU(1000, time.Minute),
					StatsEnabled: configuration.Environment == internal.Environment.DEVELOPMENT,
				})

				return &Cache{
					configuration: configuration,
					logger:        logger,
					pool:          pool,
					cache:         cache,
				}, nil
			}
		}
	}
}

func (self *Cache) Health(ctx context.Context) error {
	delay := time.NewTicker(100 * time.Millisecond)
	timeoutExceeded := time.After(300 * time.Millisecond)

	for {
		select {
		case <-timeoutExceeded:
			return ErrGeneric()
		case <-delay.C:
			err := func() error {
				if self.pool.PoolStats().TotalConns < uint32(self.configuration.CacheMinConns) {
					return ErrBelowMin()
				}

				result, err := self.pool.Ping(ctx).Result()
				if err != nil || result != "PONG" {
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

func (self *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	err := self.cache.Set(&cache.Item{
		Ctx:            ctx,
		Key:            key,
		Value:          value,
		TTL:            ttl,
		SkipLocalCache: false,
	})

	if err != nil {
		return ErrGeneric().Wrap(err)
	}

	return nil
}

type Scan struct {
	Scan func(dest interface{}) error
}

func (self *Cache) Get(ctx context.Context, key string) Scan {
	return Scan{
		Scan: func(dest interface{}) error {
			err := self.cache.Get(ctx, key, dest)
			if err != nil {
				return internalError(err)
			}

			return nil
		},
	}
}

func (self *Cache) Delete(ctx context.Context, key string) error {
	err := self.cache.Delete(ctx, key)
	if err != nil {
		return ErrGeneric().Wrap(err)
	}

	return nil
}

func (self *Cache) Close(ctx context.Context) error {
	self.logger.Info("Closing cache")

	err := self.pool.Close()
	if err != nil {
		return ErrGeneric().Wrap(err)
	}

	return nil
}

type RedisLogger struct {
	logger core.Logger
}

func NewRedisLogger(logger core.Logger) *RedisLogger {
	return &RedisLogger{
		logger: logger,
	}
}

func (self RedisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	self.logger.Infof(format, v...)
}
