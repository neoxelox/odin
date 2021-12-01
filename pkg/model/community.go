package model

import (
	"fmt"
	"time"

	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/utility"
	"github.com/rs/xid"
)

const (
	COMMUNITY_ADDRESS_MAX_LENGTH = 100
	COMMUNITY_ADDRESS_MIN_LENGTH = 1
	COMMUNITY_NAME_MAX_LENGTH    = 100
	COMMUNITY_NAME_MIN_LENGTH    = 1
)

var (
	COMMUNITY_DEFAULT_CATEGORIES = []string{
		"Suministros",
		"Desagües",
		"Cerrajería",
		"Ascensor",
		"Estructural",
		"Zonas Comunes",
		"Otros",
	}
)

type Community struct {
	class.Model
	ID         string     `db:"id"`
	Address    string     `db:"address"`
	Name       string     `db:"name"`
	Categories []string   `db:"categories"`
	PinnedIDs  []string   `db:"pinned_ids"`
	CreatedAt  time.Time  `db:"created_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}

func NewCommunity() *Community {
	now := time.Now()

	return &Community{
		ID:         xid.New().String(),
		CreatedAt:  now,
		Categories: []string{},
		PinnedIDs:  []string{},
	}
}

func (self Community) String() string {
	return fmt.Sprintf("<%s: %s>", self.Name, self.ID)
}

func (self *Community) Copy() *Community {
	return &Community{
		ID:         *utility.CopyString(&self.ID),
		Address:    *utility.CopyString(&self.Address),
		Name:       *utility.CopyString(&self.Name),
		Categories: *utility.CopyStringSlice(&self.Categories),
		PinnedIDs:  *utility.CopyStringSlice(&self.PinnedIDs),
		CreatedAt:  *utility.CopyTime(&self.CreatedAt),
		DeletedAt:  utility.CopyTime(self.DeletedAt),
	}
}
