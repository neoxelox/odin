package internal

import (
	"runtime"

	"github.com/imdario/mergo"
	"github.com/neoxelox/odin/internal/utility"
)

var Environment = struct {
	PRODUCTION  string
	DEVELOPMENT string
	TESTING     string
}{"prod", "dev", "test"}

type Configuration struct {
	Environment string

	DatabaseHost     string
	DatabasePort     int
	DatabaseSSLMode  string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseMinConns int
	DatabaseMaxConns int

	CacheHost     string
	CachePort     int
	CachePassword string
	CacheMinConns int
	CacheMaxConns int

	AppHost    string
	AppPort    int
	AppRelease string
	AppOrigins []string
	AppName    string

	TwilioBaseURL    string
	TwilioAccountSID string
	TwilioApiKey     string
	TwilioApiSecret  string
	TwilioFromPhone  string
	TwilioOriginator string
	TwilioRegion     string
	TwilioEdge       string

	SendGridApiKey    string
	SendGridFromName  string
	SendGridFromEmail string

	GracefulTimeout      int
	RequestHeaderMaxSize int
	RequestBodyMaxSize   int
	RequestFileMaxSize   int
	SessionKey           string
}

func NewConfiguration(override *Configuration) *Configuration {
	configuration := &Configuration{
		Environment: utility.GetEnvAsString("ASGARD_ENVIRONMENT", "dev"),

		DatabaseHost:     utility.GetEnvAsString("ASGARD_DATABASE_HOST", "postgres"),
		DatabasePort:     utility.GetEnvAsInt("ASGARD_DATABASE_PORT", 5432),
		DatabaseSSLMode:  utility.GetEnvAsString("ASGARD_DATABASE_SSLMODE", "disable"),
		DatabaseUser:     utility.GetEnvAsString("ODIN_DATABASE_USER", "odin"),
		DatabasePassword: utility.GetEnvAsString("ODIN_DATABASE_PASSWORD", "odin"),
		DatabaseName:     utility.GetEnvAsString("ODIN_DATABASE_NAME", "odin"),
		DatabaseMinConns: 1,
		DatabaseMaxConns: 22,

		CacheHost:     utility.GetEnvAsString("ASGARD_CACHE_HOST", "redis"),
		CachePort:     utility.GetEnvAsInt("ASGARD_CACHE_PORT", 6379),
		CachePassword: utility.GetEnvAsString("ASGARD_CACHE_PASSWORD", "redis"),
		CacheMinConns: 1,
		CacheMaxConns: 10 * runtime.GOMAXPROCS(-1),

		AppHost:    utility.GetEnvAsString("ODIN_HOST", "localhost"),
		AppPort:    utility.GetEnvAsInt("ODIN_PORT", 1111),
		AppRelease: utility.GetEnvAsString("ODIN_RELEASE", "fake"),
		AppOrigins: utility.GetEnvAsSlice("ODIN_ORIGINS", []string{"http://localhost:1111"}),
		AppName:    "odin",

		TwilioBaseURL:    utility.GetEnvAsString("ODIN_TWILIO_BASE_URL", "https://api.twilio.com/2010-04-01"),
		TwilioAccountSID: utility.GetEnvAsString("ODIN_TWILIO_ACCOUNT_SID", "fake"),
		TwilioApiKey:     utility.GetEnvAsString("ODIN_TWILIO_API_KEY", "fake"),
		TwilioApiSecret:  utility.GetEnvAsString("ODIN_TWILIO_API_SECRET", "fake"),
		TwilioFromPhone:  utility.GetEnvAsString("ODIN_TWILIO_FROM_PHONE", "fake"),
		TwilioOriginator: utility.GetEnvAsString("ODIN_TWILIO_ORIGINATOR", "Community"),
		TwilioRegion:     utility.GetEnvAsString("ODIN_TWILIO_REGION", "de1"),
		TwilioEdge:       utility.GetEnvAsString("ODIN_TWILIO_EDGE", "frankfurt"),

		SendGridApiKey:    utility.GetEnvAsString("ODIN_SENDGRID_API_KEY", "fake"),
		SendGridFromName:  utility.GetEnvAsString("ODIN_SENDGRID_FROM_NAME", "fake"),
		SendGridFromEmail: utility.GetEnvAsString("ODIN_SENDGRID_FROM_EMAIL", "fake"),

		GracefulTimeout:      15,      // 15 S
		RequestHeaderMaxSize: 1 << 10, // 1 KB
		RequestBodyMaxSize:   4 << 10, // 4 KB
		RequestFileMaxSize:   2 << 20, // 2 MB
		SessionKey:           utility.GetEnvAsString("ODIN_SESSION_KEY", "g39Ho9esU5u7ToYiaiHpDbwd5ufmoziy"),
	}

	if override == nil {
		return configuration
	}

	err := mergo.Merge(configuration, override, mergo.WithOverride)
	if err != nil {
		return configuration
	}

	return configuration
}
