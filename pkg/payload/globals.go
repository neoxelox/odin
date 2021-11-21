package payload

import (
	"net/http"

	"github.com/neoxelox/odin/internal"
)

var (
	ExcInvalidPhone = internal.NewException(http.StatusBadRequest, "ERR_INVALID_PHONE")
)
