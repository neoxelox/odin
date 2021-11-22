package pkg

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
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
	"github.com/neoxelox/odin/pkg/middleware"
	"github.com/neoxelox/odin/pkg/repository"
	"github.com/neoxelox/odin/pkg/service"
	"github.com/neoxelox/odin/pkg/usecase/auth"
	"github.com/neoxelox/odin/pkg/usecase/file"
	"github.com/neoxelox/odin/pkg/usecase/otp"
	"github.com/neoxelox/odin/pkg/usecase/session"
	"github.com/neoxelox/odin/pkg/usecase/user"
	"github.com/neoxelox/odin/pkg/view"
)

var ErrAPIGeneric = internal.NewError("API failed")

type API struct {
	class.API
}

func NewAPI(configuration internal.Configuration, logger core.Logger) (*API, error) {
	ctx := context.Background()

	apiHost := fmt.Sprintf("%s:%d", configuration.AppHost, configuration.AppPort)
	apiOrigin := fmt.Sprintf("http://%s", apiHost)
	if configuration.Environment == internal.Environment.PRODUCTION {
		apiOrigin = fmt.Sprintf("https://%s", apiHost)
	}

	/* DEPENDENCIES */

	serializer := core.NewSerializer(configuration, logger)
	validator := core.NewValidator(configuration, logger)
	renderer := core.NewRenderer(configuration, logger)
	binder := core.NewBinder(configuration, logger)
	errorHandler := core.NewErrorHandler(configuration, logger)

	database, err := database.New(ctx, 5, configuration, logger)
	if err != nil {
		return nil, ErrAPIGeneric().Wrap(err)
	}

	cache, err := cache.New(ctx, 5, configuration, logger)
	if err != nil {
		return nil, ErrAPIGeneric().Wrap(err)
	}

	server := server.New(configuration, logger, *serializer, *validator, *renderer, *binder, *errorHandler)

	api := server.Host(apiHost)

	/* MIDDLEWARES */

	httpsRedirectMiddleware := echoMiddleware.HTTPSRedirect()
	removeTrailingSlashMiddleware := echoMiddleware.RemoveTrailingSlash()
	loggerMiddleware := internalMiddleware.NewLoggerMiddleware(configuration, logger).Handle
	recoverMiddleware := echoMiddleware.Recover()
	corsMiddleware := echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{apiOrigin},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
		AllowHeaders: []string{"*"},
		MaxAge:       86400,
	})
	bodyLimitMiddleware := echoMiddleware.BodyLimitWithConfig(echoMiddleware.BodyLimitConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Path(), strings.TrimPrefix(internal.ASSETS_PATH, "."))
		},
		Limit: utility.SizeToString(configuration.RequestBodyMaxSize),
	})
	fileLimitMiddleware := echoMiddleware.BodyLimitWithConfig(echoMiddleware.BodyLimitConfig{
		Limit: utility.SizeToString(configuration.RequestFileMaxSize),
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
			return strings.HasPrefix(c.Path(), strings.TrimPrefix(internal.ASSETS_PATH, "."))
		},
		Level: 5,
	})
	timeoutMiddleware := echoMiddleware.TimeoutWithConfig(echoMiddleware.TimeoutConfig{
		Timeout:                    time.Duration(configuration.GracefulTimeout*2/3) * time.Second,
		ErrorMessage:               internal.ExcRequestTimeout.JSON(),
		OnTimeoutRouteErrorHandler: errorHandler.Handle,
	})

	// TODO: Fix bug timeoutMiddleware calls two times errorHandler.Handle if not timeout
	// TODO: CSRF Middleware
	// TODO: Rate Limiter Middleware
	// TODO: Request ID Middleware (open telemetry)

	/* SERVICES */

	smsService := service.NewSMSService(configuration, logger)
	emailService := service.NewEmailService(configuration, logger)

	/* REPOSITORIES  */

	_ = repository.NewCommunityRepository(configuration, logger, *database)
	_ = repository.NewInvitationRepository(configuration, logger, *database)
	_ = repository.NewMembershipRepository(configuration, logger, *database)
	otpRepository := repository.NewOTPRepository(configuration, logger, *database)
	_ = repository.NewPostRepository(configuration, logger, *database)
	sessionRepository := repository.NewSessionRepository(configuration, logger, *database)
	userRepository := repository.NewUserRepository(configuration, logger, *database)

	/* USECASES */

	fileCreator := file.NewCreatorUsecase(configuration, logger)
	fileGetter := file.NewGetterUsecase(configuration, logger)

	otpCreator := otp.NewCreatorUsecase(configuration, logger, *database, *otpRepository, *smsService, *emailService)
	otpVerifier := otp.NewVerifierUsecase(configuration, logger, *otpRepository)

	sessionCreator := session.NewCreatorUsecase(configuration, logger, *database, *sessionRepository, *userRepository)

	userCreator := user.NewCreatorUsecase(configuration, logger, *userRepository)
	userUpdater := user.NewUpdaterUsecase(configuration, logger, *userRepository)

	authCreator := auth.NewCreatorUsecase(configuration, logger)
	authVerifier := auth.NewVerifierUsecase(configuration, logger, *sessionRepository, *userRepository)
	authLogger := auth.NewLoggerUsecase(configuration, logger, *database, *otpVerifier, *userCreator,
		*sessionCreator, *authCreator, *otpRepository, *userRepository)

	/* VIEWS */

	healthView := internalView.NewHealthView(configuration, logger, *database, *cache)
	fileView := view.NewFileView(configuration, logger, *fileCreator, *fileGetter)
	authView := view.NewAuthView(configuration, logger, *otpCreator, *authLogger)
	userView := view.NewUserView(configuration, logger, *userUpdater)

	/* MIDDLEWARES */

	authMiddleware := middleware.NewAuthMiddleware(configuration, logger, *authVerifier).Handle

	/* ROUTES */

	if configuration.Environment == internal.Environment.PRODUCTION {
		server.Use(httpsRedirectMiddleware)
	}
	server.Use(removeTrailingSlashMiddleware)

	api.Use(loggerMiddleware)
	api.Use(recoverMiddleware)
	api.Use(corsMiddleware)
	api.Use(bodyLimitMiddleware)
	api.Use(secureMiddleware)
	api.Use(gzipMiddleware)
	api.Use(timeoutMiddleware)

	api.GET("/health", healthView.Handle(healthView.GetHealth()))

	// NOT AUTHENTICATED

	api.POST("/login/start", authView.Handle(authView.LoginStart()))
	api.POST("/login/end", authView.Handle(authView.LoginEnd()))

	apiV1 := api.Group("/v1")

	// AUTHENTICATED

	apiV1 = apiV1.Group("", authMiddleware)

	apiV1.POST("/file", fileView.Handle(fileView.Post()), fileLimitMiddleware)
	apiV1.GET("/file/:name", fileView.Handle(fileView.Get()), fileLimitMiddleware)

	apiV1.POST("/user/profile", userView.Handle(userView.PostProfile()))

	return &API{
		API: *class.NewAPI(
			// ON START
			func() error {
				return nil
			},
			// ON CLOSE
			func(ctx context.Context) error {
				var err error
				if errC := smsService.Close(ctx); errC != nil {
					err = errors.CombineErrors(err, errC)
				}

				return err
			}, configuration, logger, *database, *cache, *server),
	}, nil
}
