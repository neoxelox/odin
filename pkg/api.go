package pkg

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/cache"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	internalMiddleware "github.com/neoxelox/odin/internal/middleware"
	"github.com/neoxelox/odin/internal/server"
	"github.com/neoxelox/odin/internal/utility"
	internalView "github.com/neoxelox/odin/internal/view"
	"github.com/neoxelox/odin/pkg/usecase/file"
	"github.com/neoxelox/odin/pkg/view"
)

var ErrAPIGeneric = internal.NewError("API failed")

type API struct {
	class.API
}

func NewAPI(configuration internal.Configuration, logger core.Logger) (*API, error) {
	ctx := context.Background()

	database, err := database.New(ctx, 5, configuration, logger)
	if err != nil {
		return nil, ErrAPIGeneric().Wrap(err)
	}

	cache, err := cache.New(ctx, 5, configuration, logger)
	if err != nil {
		return nil, ErrAPIGeneric().Wrap(err)
	}

	serializer := core.NewSerializer(configuration, logger)
	validator := core.NewValidator(configuration, logger)
	renderer := core.NewRenderer(configuration, logger)
	binder := core.NewBinder(configuration, logger)
	errorHandler := core.NewErrorHandler(configuration, logger)

	server := server.New(configuration, logger, *serializer, *validator, *renderer, *binder, *errorHandler)

	if configuration.Environment == internal.Environment.PRODUCTION {
		server.Use(echoMiddleware.HTTPSRedirect())
	}
	server.Use(echoMiddleware.RemoveTrailingSlash())

	apiHost := fmt.Sprintf("%s:%d", configuration.AppHost, configuration.AppPort)
	apiOrigin := fmt.Sprintf("http://%s", apiHost)
	if configuration.Environment == internal.Environment.PRODUCTION {
		apiOrigin = fmt.Sprintf("https://%s", apiHost)
	}

	api := server.Host(apiHost)

	loggerMiddleware := internalMiddleware.NewLoggerMiddleware(configuration, logger)
	recoverMiddleware := echoMiddleware.Recover()
	corsMiddleware := echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{apiOrigin},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
		AllowHeaders: []string{"*"},
		MaxAge:       86400,
	})
	bodyLimitMiddleware := echoMiddleware.BodyLimitWithConfig(echoMiddleware.BodyLimitConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "/static") || strings.Contains(c.Path(), "/file")
		},
		Limit: utility.SizeToString(configuration.RequestBodyMaxSize),
	})
	secureMiddleware := echoMiddleware.SecureWithConfig(echoMiddleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: false,
		HSTSPreloadEnabled:    true,
		ContentSecurityPolicy: fmt.Sprintf("default-src %s", apiOrigin),
		CSPReportOnly:         false,
		ReferrerPolicy:        "same-origin",
	})
	gzipMiddleware := echoMiddleware.GzipWithConfig(echoMiddleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "/static") || strings.Contains(c.Path(), "/file")
		},
		Level: 5,
	})
	timeoutMiddleware := echoMiddleware.TimeoutWithConfig(echoMiddleware.TimeoutConfig{
		Timeout:                    time.Duration(configuration.GracefulTimeout*2/3) * time.Second,
		ErrorMessage:               internal.ErrRequestTimeout.JSON(),
		OnTimeoutRouteErrorHandler: errorHandler.Handle,
	}) // TODO: Fix bug timeoutMiddleware calls two times errorHandler.Handle if not timeout

	api.Use(loggerMiddleware.Handle)
	api.Use(recoverMiddleware)
	api.Use(corsMiddleware)
	api.Use(bodyLimitMiddleware)
	api.Use(secureMiddleware)
	api.Use(gzipMiddleware)
	api.Use(timeoutMiddleware)
	// TODO: CSRF Middleware
	// TODO: Rate Limiter Middleware
	// TODO: Request ID Middleware (open telemetry)

	healthView := internalView.NewHealthView(configuration, logger, *database, *cache)
	api.GET("/health", healthView.Handle(healthView.GetHealth()))

	apiV1 := api.Group("/v1")

	fileLimitMiddleware := echoMiddleware.BodyLimitWithConfig(echoMiddleware.BodyLimitConfig{
		Limit: utility.SizeToString(configuration.RequestFileMaxSize),
	})

	fileCreator := file.NewCreatorUsecase(configuration, logger)
	fileGetter := file.NewGetterUsecase(configuration, logger)
	fileView := view.NewFileView(configuration, logger, *fileCreator, *fileGetter)
	apiV1.POST("/file", fileView.Handle(fileView.Post()), fileLimitMiddleware)
	apiV1.GET("/file/:file_name", fileView.Handle(fileView.Get()), fileLimitMiddleware)

	loginView := view.NewLoginView(configuration, logger)
	apiV1.POST("/login/start", loginView.Handle(loginView.PostStart()))

	_ = api.Group("") // auth required

	return &API{
		API: *class.NewAPI(
			func() error {
				return nil
			}, func(ctx context.Context) error {
				return nil
			}, configuration, logger, *database, *cache, *server),
	}, nil
}
