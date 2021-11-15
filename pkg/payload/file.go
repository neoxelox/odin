package payload

import (
	"github.com/neoxelox/odin/internal/class"
)

type PostFileResponse struct {
	class.Payload
	FileURL string `json:"file_url"`
}

type GetFileRequest struct {
	class.Payload
	FileName string `param:"file_name" validate:"required"`
}
