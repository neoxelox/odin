package view

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/cache"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
)

type HealthView struct {
	class.View
	database database.Database
	cache    cache.Cache
}

func NewHealthView(configuration internal.Configuration, logger core.Logger, database database.Database, cache cache.Cache) *HealthView {
	return &HealthView{
		View:     *class.NewView(configuration, logger),
		database: database,
		cache:    cache,
	}
}

func (self *HealthView) GetHealth(ctx echo.Context) error {
	return self.Handle(ctx, class.Endpoint{}, func() error {
		err := self.database.Health(ctx.Request().Context())
		if err != nil {
			return internal.ExcServerUnavailable.Cause(err)
		}

		err = self.cache.Health(ctx.Request().Context())
		if err != nil {
			return internal.ExcServerUnavailable.Cause(err)
		}

		return ctx.String(http.StatusOK, "OK\n")
	})
}
