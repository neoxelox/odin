package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/utility"
	"github.com/rs/xid"
)

const (
	POST_MESSAGE_MAX_LENGTH = 280
	POST_MESSAGE_MIN_LENGTH = 1
)

type Post struct {
	class.Model
	ID            string    `db:"id"`
	ThreadID      *string   `db:"thread_id"`
	CreatorID     string    `db:"creator_id"`
	LastHistoryID *string   `db:"last_history_id"`
	Type          string    `db:"type"`
	Priority      *int      `db:"priority"`
	RecipientIDs  *[]string `db:"recipient_ids"`
	VoterIDs      []string  `db:"voter_ids"`
	CreatedAt     time.Time `db:"created_at"`
}

var PostType = struct {
	PUBLICATION string
	ISSUE       string
	EVENT       string
	Has         func(typee string) bool
}{"PUBLICATION", "ISSUE", "EVENT", func(typee string) bool {
	return typee == "PUBLICATION" || typee == "ISSUE" || typee == "EVENT"
}}

var PostState = struct {
	PENDING     string
	IN_PROGRESS string
	REJECTED    string
	ACCEPTED    string
	RESOLVED    string
	Has         func(role string) bool
}{"PENDING", "IN_PROGRESS", "REJECTED", "ACCEPTED", "RESOLVED",
	func(role string) bool {
		return role == "PENDING" || role == "IN_PROGRESS" ||
			role == "REJECTED" || role == "ACCEPTED" || role == "RESOLVED"
	},
}

type PostWidgets struct {
	Poll *map[string][]string
}

type PostMedia struct {
	Pictures *[]string
	Videos   *[]string
	Audios   *[]string
}

func NewPost() *Post {
	now := time.Now()

	return &Post{
		ID:        xid.New().String(),
		CreatedAt: now,
	}
}

func (self Post) String() string {
	return fmt.Sprintf("<%s: %s>", self.Type, self.ID)
}

func (self *Post) Copy() *Post {
	return &Post{
		ID:            *utility.CopyString(&self.ID),
		ThreadID:      utility.CopyString(self.ThreadID),
		CreatorID:     *utility.CopyString(&self.CreatorID),
		LastHistoryID: utility.CopyString(self.LastHistoryID),
		Type:          *utility.CopyString(&self.Type),
		Priority:      utility.CopyInt(self.Priority),
		RecipientIDs:  utility.CopyStringSlice(self.RecipientIDs),
		VoterIDs:      *utility.CopyStringSlice(&self.VoterIDs),
		CreatedAt:     *utility.CopyTime(&self.CreatedAt),
	}
}

type PostHistory struct {
	class.Model
	ID         string      `db:"id"`
	PostID     string      `db:"post_id"`
	Message    string      `db:"message"`
	Categories []string    `db:"categories"`
	State      *string     `db:"state"`
	Widgets    PostWidgets `db:"widgets"`
	Media      PostMedia   `db:"media"`
	CreatedAt  time.Time   `db:"created_at"`
}

func NewPostHistory() *PostHistory {
	now := time.Now()

	return &PostHistory{
		ID:         xid.New().String(),
		CreatedAt:  now,
		Categories: []string{},
	}
}

func (self PostHistory) String() string {
	return fmt.Sprintf("<%s: %s>", self.Message, self.ID)
}

func (self *PostHistory) Copy() *PostHistory {
	return &PostHistory{
		ID:         *utility.CopyString(&self.ID),
		PostID:     *utility.CopyString(&self.PostID),
		Message:    *utility.CopyString(&self.Message),
		Categories: *utility.CopyStringSlice(&self.Categories),
		State:      utility.CopyString(self.State),
		Widgets: PostWidgets{
			Poll: utility.CopyStringSliceMap(self.Widgets.Poll),
		},
		Media: PostMedia{
			Pictures: utility.CopyStringSlice(self.Media.Pictures),
			Videos:   utility.CopyStringSlice(self.Media.Videos),
			Audios:   utility.CopyStringSlice(self.Media.Audios),
		},
		CreatedAt: *utility.CopyTime(&self.CreatedAt),
	}
}
