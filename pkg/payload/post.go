package payload

import (
	"time"

	"github.com/neoxelox/odin/internal/class"
)

type Post struct {
	ID           string    `json:"id"`
	ThreadID     *string   `json:"thread_id"`
	CreatorID    string    `json:"creator_id"`
	Type         string    `json:"type"`
	Priority     *int      `json:"priority"`
	RecipientIDs *[]string `json:"recipient_ids"`
	VoterIDs     []string  `json:"voter_ids"`
	CreatedAt    time.Time `json:"created_at"`
	PostHistory
}

type PostWidgets struct {
	Poll *map[string][]string `json:"poll,omitempty"`
}

type PostHistory struct {
	Message    string      `json:"message"`
	Categories []string    `json:"categories"`
	State      *string     `json:"state"`
	Media      []string    `json:"media"`
	Widgets    PostWidgets `json:"widgets"`
	CreatedAt  time.Time   `json:"created_at,omitempty"`
}

type PostPostRequest struct {
	class.Payload
	CommunityID  string    `param:"id" validate:"required"`
	Type         string    `json:"type" validate:"required"`
	ThreadID     *string   `json:"thread_id" validate:"omitempty,required"`
	Priority     *int      `json:"priority" validate:"omitempty,required"`
	RecipientIDs *[]string `json:"recipient_ids" validate:"omitempty,required"`
	Message      string    `json:"message" validate:"required"`
	Categories   *[]string `json:"categories" validate:"omitempty,required"`
	State        *string   `json:"state" validate:"omitempty,required"`
	Media        *[]string `json:"media" validate:"omitempty,required"`
	Widgets      *struct {
		PollOptions *[]string `json:"poll_options" validate:"omitempty,required"`
	} `json:"widgets" validate:"omitempty,required"`
}

type PostPostResponse struct {
	class.Payload
	Post
}
