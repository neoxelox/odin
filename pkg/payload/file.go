package payload

import (
	"github.com/neoxelox/odin/internal/class"
)

type PostFileResponse struct {
	class.Payload
	URL string `json:"url"`
}

type GetFileRequest struct {
	class.Payload
	Name string `param:"name" validate:"required"`
}
