package payload

import (
	"net/http"

	"github.com/neoxelox/odin/internal"
)

var (
	ErrInvalidPhone = internal.NewException(http.StatusBadRequest, "ERR_INVALID_PHONE")
)
