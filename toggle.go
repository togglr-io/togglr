package toggle

import (
	"context"
	"time"

	"github.com/eriktate/toggle/rules"
	"github.com/eriktate/toggle/uid"
)

// Metadata is included with the initial request for Toggles when a new client initializes. It's used to
// evaluate rules and determine the final value for each toggle.
type Metadata struct {
}

// A Toggle represents a key and the set of rules that determine the value that should be returned for it.
type Toggle struct {
	ID          uid.UID      `json:"id"`
	Key         string       `json:"key"`
	Description string       `json:"description"`
	Rules       []rules.Rule `json:"rules"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

// ListTogglesReq defines the search parameters that will be used when generating a list of toggles.
type ListTogglesReq struct{}

// A ToggleService performs basic CRUD operations on toggles.
type ToggleService interface {
	CreateToggle(ctx context.Context, toggle Toggle) (uid.UID, error)
	FetchToggle(ctx context.Context, id uid.UID) (Toggle, error)
	ListToggles(ctx context.Context, req ListTogglesReq) ([]Toggle, error)
}
