package togglr

import (
	"context"
	"time"

	"github.com/togglr-io/togglr/rules"
	"github.com/togglr-io/togglr/uid"
)

// Metadata is included with the initial request for Toggles when a new client initializes. It's used to
// evaluate rules and determine the final value for each toggle.
type Metadata struct {
}

// An ID is a wrapper struct for working with JSON payloads containing
// an ID. This is used to differentiate between creates and updates
// as well as a return type for certain operations.
type ID struct {
	ID uid.UID `json:"id"`
}

// An Account represents a grouping of Users and Toggles.
type Account struct {
	ID        uid.UID   `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// A Toggle represents a key and the set of rules that determine the value that should be returned for it.
type Toggle struct {
	ID          uid.UID     `json:"id" db:"id"`
	AccountID   uid.UID     `json:"accountId" db:"account_id"`
	Key         string      `json:"key" db:"key"`
	Description string      `json:"description" db:"description"`
	Active      bool        `json:"active" db:"active"`
	Rules       rules.Rules `json:"rules" db:"rules"`
	CreatedAt   time.Time   `json:"createdAt" db:"created_at" goqu:"skipinsert"`
	UpdatedAt   time.Time   `json:"updatedAt" db:"updated_at" goqu:"skipinsert"`
}

// An UpdateToggleReq contains all of the fields that are possible to update on a Toggle. The main difference from
// the Toggle struct is that some of the fields are pointers to differentiate from a field being omitted and an
// actual update containing the zero value
type UpdateToggleReq struct {
	ID          uid.UID      `json:"id" db:"id"`
	Key         *string      `json:"key,omitempty" db:"key"`
	Description *string      `json:"description,omitempty" db:"description"`
	Active      *bool        `json:"active" db:"active"`
	Rules       []rules.Rule `json:"rules" db:"-"`
}

// ListTogglesReq defines the search parameters that will be used when generating a list of toggles.
type ListTogglesReq struct {
}

// A ToggleService performs basic CRUD operations on toggles.
type ToggleService interface {
	CreateToggle(ctx context.Context, toggle Toggle) (uid.UID, error)
	UpdateToggle(ctx context.Context, req UpdateToggleReq) error
	FetchToggle(ctx context.Context, id uid.UID) (Toggle, error)
	ListToggles(ctx context.Context, req ListTogglesReq) ([]Toggle, error)
	DeleteToggle(ctx context.Context, id uid.UID) error
}
