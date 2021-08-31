package pg

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/togglr-io/togglr"
	"github.com/togglr-io/togglr/uid"
)

func (c Client) FetchKeys(ctx context.Context, accountID uid.UID) ([]togglr.MetadataKey, error) {
	keys := []togglr.MetadataKey{}
	query := c.db.
		From("metadata_keys").
		Where(goqu.Ex{"account_id": accountID})

	if err := query.ScanStructs(&keys); err != nil {
		return nil, err
	}

	return keys, nil
}

func (c Client) PushKeys(ctx context.Context, accountID uid.UID, keys ...string) error {
	newKeys := make([]togglr.MetadataKey, len(keys))
	for idx, key := range keys {
		newKeys[idx] = togglr.MetadataKey{
			ID:        uid.New(),
			AccountID: accountID,
			Key:       key,
		}
	}

	if _, err := c.db.Insert("metadata_keys").Rows(newKeys).OnConflict(goqu.DoNothing()).Executor().Exec(); err != nil {
		return err
	}

	return nil
}
