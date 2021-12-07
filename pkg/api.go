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
	"github.com/scylladb/go-set/strset"

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
	"github.com/neoxelox/odin/pkg/usecase/community"
	"github.com/neoxelox/odin/pkg/usecase/file"
	"github.com/neoxelox/odin/pkg/usecase/invitation"
	"github.com/neoxelox/odin/pkg/usecase/otp"
	"github.com/neoxelox/odin/pkg/usecase/post"
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

	apiOrigin := fmt.Sprintf("http://%s", configuration.AppHost)
	if configuration.Environment == internal.Environment.PRODUCTION {
		apiOrigin = fmt.Sprintf("https://%s", configuration.AppHost)
	}
	apiOrigins := strset.New(append([]string{apiOrigin}, configuration.AppOrigins...)...).List()

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

	api := server.Host(configuration.AppHost)

	/* MIDDLEWARES */

	// httpsRedirectMiddleware := echoMiddleware.HTTPSRedirect()
	removeTrailingSlashMiddleware := echoMiddleware.RemoveTrailingSlash()
	loggerMiddleware := internalMiddleware.NewLoggerMiddleware(configuration, logger).Handle
	recoverMiddleware := echoMiddleware.Recover()
	corsMiddleware := echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: apiOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodHead, http.MethodOptions},
		AllowHeaders: []string{"*"},
		MaxAge:       86400,
	})
	bodyLimitMiddleware := echoMiddleware.BodyLimitWithConfig(echoMiddleware.BodyLimitConfig{
		Skipper: func(c echo.Context) bool {
			return c.Request().Method == http.MethodPost && c.Path() == "/file"
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
		ContentSecurityPolicy: fmt.Sprintf("default-src %s", strings.Join(apiOrigins, " ")),
		CSPReportOnly:         false,
		ReferrerPolicy:        "same-origin",
	})
	gzipMiddleware := echoMiddleware.GzipWithConfig(echoMiddleware.GzipConfig{
		Skipper: func(c echo.Context) bool { // TODO: Refactorize this!
			return strings.HasPrefix(c.Path(), "/file") || strings.HasPrefix(c.Path(), strings.TrimPrefix(internal.ASSETS_PATH, "."))
		},
		Level: 6,
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

	communityRepository := repository.NewCommunityRepository(configuration, logger, *database)
	invitationRepository := repository.NewInvitationRepository(configuration, logger, *database)
	membershipRepository := repository.NewMembershipRepository(configuration, logger, *database)
	otpRepository := repository.NewOTPRepository(configuration, logger, *database)
	postRepository := repository.NewPostRepository(configuration, logger, *database)
	sessionRepository := repository.NewSessionRepository(configuration, logger, *database)
	userRepository := repository.NewUserRepository(configuration, logger, *database)

	/* USECASES */

	fileCreator := file.NewCreatorUsecase(configuration, logger)
	fileGetter := file.NewGetterUsecase(configuration, logger)

	otpCreator := otp.NewCreatorUsecase(configuration, logger, *database, *renderer, *otpRepository, *smsService, *emailService)
	otpVerifier := otp.NewVerifierUsecase(configuration, logger, *otpRepository)

	sessionCreator := session.NewCreatorUsecase(configuration, logger, *database, *sessionRepository, *userRepository)

	userGetter := user.NewGetterUsecase(configuration, logger, *userRepository)
	userCreator := user.NewCreatorUsecase(configuration, logger, *userRepository)
	userUpdater := user.NewUpdaterUsecase(configuration, logger, *database, *userRepository, *otpRepository, *otpVerifier)
	userDeleter := user.NewDeleterUsecase(configuration, logger, *userRepository)

	communityGetter := community.NewGetterUsecase(configuration, logger, *communityRepository, *membershipRepository, *userRepository, *invitationRepository)
	communityJoiner := community.NewJoinerUsecase(configuration, logger, *communityRepository, *membershipRepository)
	communityLeaver := community.NewLeaverUsecase(configuration, logger, *communityRepository, *membershipRepository)
	communityCreator := community.NewCreatorUsecase(configuration, logger, *database, *communityJoiner, *communityRepository)

	invitationGetter := invitation.NewGetterUsecase(configuration, logger, *invitationRepository)
	invitationAccepter := invitation.NewAccepterUsecase(configuration, logger, *database, *communityJoiner, *invitationRepository)
	invitationRejecter := invitation.NewRejecterUsecase(configuration, logger, *invitationRepository)
	invitationCreator := invitation.NewCreatorUsecase(configuration, logger, *database, *invitationRepository, *membershipRepository, *userRepository)

	postGetter := post.NewGetterUsecase(configuration, logger, *postRepository, *membershipRepository)
	postCreator := post.NewCreatorUsecase(configuration, logger, *database, *postRepository, *membershipRepository)
	postUpdater := post.NewUpdaterUsecase(configuration, logger, *database, *postRepository, *membershipRepository)
	postVoter := post.NewVoterUsecase(configuration, logger, *postRepository, *membershipRepository)
	postUnvoter := post.NewUnvoterUsecase(configuration, logger, *postRepository, *membershipRepository)
	postPollVoter := post.NewPollVoterUsecase(configuration, logger, *postRepository, *membershipRepository)
	postPinner := post.NewPinnerUsecase(configuration, logger, *postRepository, *membershipRepository, *communityRepository)
	postUnpinner := post.NewUnpinnerUsecase(configuration, logger, *postRepository, *membershipRepository, *communityRepository)

	authCreator := auth.NewCreatorUsecase(configuration, logger)
	authVerifier := auth.NewVerifierUsecase(configuration, logger, *sessionRepository, *userRepository)
	authLogger := auth.NewLoggerUsecase(configuration, logger, *database, *otpVerifier, *userCreator,
		*sessionCreator, *authCreator, *otpRepository, *userRepository, *sessionRepository, *invitationCreator)

	/* VIEWS */

	healthView := internalView.NewHealthView(configuration, logger, *database, *cache)
	fileView := view.NewFileView(configuration, logger, *fileCreator, *fileGetter)
	authView := view.NewAuthView(configuration, logger, *otpCreator, *authLogger)
	userView := view.NewUserView(configuration, logger, *userGetter, *userUpdater, *userDeleter, *otpCreator)
	communityView := view.NewCommunityView(configuration, logger, *communityCreator, *communityGetter, *communityLeaver, *invitationCreator)
	invitationView := view.NewInvitationView(configuration, logger, *invitationGetter, *invitationAccepter, *invitationRejecter)
	postView := view.NewPostView(configuration, logger, *postGetter, *postCreator, *postUpdater, *postVoter, *postUnvoter, *postPollVoter, *postPinner, *postUnpinner, *postRepository)

	/* MIDDLEWARES */

	authMiddleware := middleware.NewAuthMiddleware(configuration, logger, *authVerifier).Handle

	/* ROUTES */

	// if configuration.Environment == internal.Environment.PRODUCTION {
	// 	server.Use(httpsRedirectMiddleware)
	// }
	server.Use(removeTrailingSlashMiddleware)

	api.Use(loggerMiddleware)
	api.Use(recoverMiddleware)
	api.Use(corsMiddleware)
	api.Use(bodyLimitMiddleware)
	api.Use(secureMiddleware)
	api.Use(gzipMiddleware)
	api.Use(timeoutMiddleware)

	// NOT AUTHENTICATED

	api.GET("/health", healthView.GetHealth)

	api.POST("/login/start", authView.PostLoginStart)
	api.POST("/login/end", authView.PostLoginEnd)

	// TODO: Refactorize this, discuss whether files and assets should be merged
	api.GET("/file/:name", fileView.GetFile)

	// VERSIONED

	apiV1 := api.Group("/v1")

	// AUTHENTICATED

	api = api.Group("", authMiddleware)

	api.POST("/file", fileView.PostFile, fileLimitMiddleware)
	api.POST("/logout", authView.PostLogout)

	// VERSIONED

	apiV1 = api.Group("/v1")

	apiV1.GET("/user/profile", userView.GetProfile)
	apiV1.POST("/user/profile", userView.PostProfile)
	apiV1.POST("/user/email/start", userView.PostEmailStart)
	apiV1.POST("/user/email/end", userView.PostEmailEnd)
	apiV1.POST("/user/phone/start", userView.PostPhoneStart)
	apiV1.POST("/user/phone/end", userView.PostPhoneEnd)
	apiV1.DELETE("/user", userView.DeleteUser)

	apiV1.GET("/community", communityView.GetCommunityList)
	apiV1.GET("/community/:id", communityView.GetCommunity)
	apiV1.POST("/community", communityView.PostCommunity)
	apiV1.POST("/community/:id/invite", communityView.PostCommunityInvite)
	apiV1.POST("/community/:id/leave", communityView.PostCommunityLeave)

	apiV1.GET("/invitation", invitationView.GetInvitationList)
	apiV1.POST("/invitation/:id/accept", invitationView.PostInvitationAccept)
	apiV1.POST("/invitation/:id/reject", invitationView.PostInvitationReject)

	apiV1.GET("/community/:id/user", communityView.GetCommunityUserList)
	apiV1.GET("/community/:community_id/user/:membership_id", communityView.GetCommunityUser)

	apiV1.GET("/community/:id/post", postView.GetPostList)
	apiV1.GET("/community/:community_id/post/:post_id", postView.GetPost)
	apiV1.GET("/community/:community_id/post/:post_id/history", postView.GetPostHistory)
	apiV1.GET("/community/:community_id/post/:post_id/thread", postView.GetPostThread)
	apiV1.POST("/community/:id/post", postView.PostPost)
	apiV1.PUT("/community/:community_id/post/:post_id", postView.PutPost)
	apiV1.POST("/community/:community_id/post/:post_id/vote", postView.PostVotePost)
	apiV1.POST("/community/:community_id/post/:post_id/unvote", postView.PostUnvotePost)
	apiV1.POST("/community/:community_id/post/:post_id/widget/poll/vote", postView.PostVotePostPoll)
	apiV1.POST("/community/:community_id/post/:post_id/pin", postView.PostPinPost)
	apiV1.POST("/community/:community_id/post/:post_id/unpin", postView.PostUnpinPost)

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
