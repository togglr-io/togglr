package pg

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/togglr-io/togglr"
	"github.com/togglr-io/togglr/uid"
)

// CreateToggle creates a new Toggle in postgres. If the toggle doen't already have an ID, one will be
// generated
func (c Client) CreateToggle(ctx context.Context, toggle togglr.Toggle) (uid.UID, error) {
	// if no ID is provided, generate one
	if !toggle.ID.IsNull() {
		toggle.ID = uid.New()
	}

	if _, err := c.db.Insert("toggles").Rows(toggle).Executor().ExecContext(ctx); err != nil {
		return toggle.ID, err
	}

	return toggle.ID, nil
}

// UpdateToggle updates an existing Toggle in postgres
func (c Client) UpdateToggle(ctx context.Context, toggle togglr.Toggle) error {
	// TODO (etate): This is a super naive update. Should probably be a bit more perscriptive.
	if _, err := c.db.Update("toggles").Set(toggle).Executor().Exec(); err != nil {
		return err
	}

	return nil
}

// FetchToggle queries a single Toggle from postgres
func (c Client) FetchToggle(ctx context.Context, id uid.UID) (togglr.Toggle, error) {
	var tog togglr.Toggle
	ds := c.db.From("toggles").Where(goqu.Ex{"id": id})
	if _, err := ds.ScanStruct(&tog); err != nil {
		return tog, err
	}

	return tog, nil
}

// ListToggles queries a slice of Toggles from postgres
func (c Client) ListToggles(ctx context.Context, req togglr.ListTogglesReq) ([]togglr.Toggle, error) {
	// default to instantiated value so that we return an empty slice instead of null when there's no results
	toggles := []togglr.Toggle{}
	if err := c.db.From("toggles").ScanStructs(&toggles); err != nil {
		return nil, err
	}

	return toggles, nil
}

// DeleteToggle deletes a Toggle from postgres
func (c Client) DeleteToggle(ctx context.Context, id uid.UID) error {
	del := c.db.Delete("toggles").Where(goqu.Ex{"id": id}).Executor()
	if _, err := del.Exec(); err != nil {
		return err
	}

	return nil
}
