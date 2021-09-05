package mock

import (
	"context"

	"github.com/togglr-io/togglr"
	"github.com/togglr-io/togglr/uid"
)

type MetadataService struct {
	FetchKeysFn     func(ctx context.Context, accountID uid.UID) ([]togglr.MetadataKey, error)
	FetchKeysCalled int

	PushKeysFn     func(ctx context.Context, accountID uid.UID, keys ...string) error
	PushKeysCalled int

	Error error
}

func NewMetadataService(err error) *MetadataService {
	return &MetadataService{Error: err}
}

func (m *MetadataService) FetchKeys(ctx context.Context, accountID uid.UID) ([]togglr.MetadataKey, error) {
	m.FetchKeysCalled++
	if m.FetchKeysFn != nil {
		return m.FetchKeysFn(ctx, accountID)
	}

	return nil, m.Error
}

func (m *MetadataService) PushKeys(ctx context.Context, accountID uid.UID, keys ...string) error {
	m.PushKeysCalled++
	if m.PushKeysFn != nil {
		return m.PushKeysFn(ctx, accountID, keys...)
	}

	return m.Error
}
