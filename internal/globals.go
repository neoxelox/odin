package internal

import "net/http"

type ContextKey string

const (
	ASSETS_PATH  = "./assets"
	FILES_PATH   = ASSETS_PATH + "/files"
	IMAGES_PATH  = ASSETS_PATH + "/images"
	SCRIPTS_PATH = ASSETS_PATH + "/scripts"
	STYLES_PATH  = ASSETS_PATH + "/styles"

	MIGRATIONS_PATH = "./migrations"

	TEMPLATES_PATH = "./templates"
)

var (
	// ExcServerGeneric generic server exception.
	ExcServerGeneric = NewException(http.StatusInternalServerError, "ERR_SERVER_GENERIC")

	// ExcServerUnavailable server Unavailable exception.
	ExcServerUnavailable = NewException(http.StatusServiceUnavailable, "ERR_SERVER_UNAVAILABLE")

	// ExcRequestTimeout request timeout exception.
	ExcRequestTimeout = NewException(http.StatusRequestTimeout, "ERR_REQUEST_TIMEOUT")

	// ExcClientGeneric generic client exception.
	ExcClientGeneric = NewException(http.StatusBadRequest, "ERR_CLIENT_GENERIC")

	// ExcInvalidRequest invalid request exception.
	ExcInvalidRequest = NewException(http.StatusBadRequest, "ERR_INVALID_REQUEST")

	// ExcInvalidRequest not found exception.
	ExcNotFound = NewException(http.StatusNotFound, "ERR_NOT_FOUND")

	// ExcUnauthorized unauthorized exception.
	ExcUnauthorized = NewException(http.StatusUnauthorized, "ERR_UNAUTHORIZED")
)
