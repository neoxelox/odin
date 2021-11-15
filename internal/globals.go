package internal

import "net/http"

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
	// ErrServerGeneric generic server exception.
	ErrServerGeneric = NewException(http.StatusInternalServerError, "ERR_SERVER_GENERIC")

	// ErrServerUnavailable server Unavailable exception.
	ErrServerUnavailable = NewException(http.StatusServiceUnavailable, "ERR_SERVER_UNAVAILABLE")

	// ErrRequestTimeout request timeout exception.
	ErrRequestTimeout = NewException(http.StatusRequestTimeout, "ERR_REQUEST_TIMEOUT")

	// ErrClientGeneric generic client exception.
	ErrClientGeneric = NewException(http.StatusBadRequest, "ERR_CLIENT_GENERIC")

	// ErrInvalidRequest invalid request exception.
	ErrInvalidRequest = NewException(http.StatusBadRequest, "ERR_INVALID_REQUEST")

	// ErrInvalidRequest not found exception.
	ErrNotFound = NewException(http.StatusNotFound, "ERR_NOT_FOUND")
)
