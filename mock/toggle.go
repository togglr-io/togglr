package mock

import (
	"context"

	"github.com/togglr-io/togglr"
	"github.com/togglr-io/togglr/uid"
)

type ToggleService struct {
	CreateToggleFn     func(ctx context.Context, toggle togglr.Toggle) (uid.UID, error)
	CreateToggleCalled int

	UpdateToggleFn     func(ctx context.Context, req togglr.UpdateToggleReq) error
	UpdateToggleCalled int

	FetchToggleFn     func(ctx context.Context, id uid.UID) (togglr.Toggle, error)
	FetchToggleCalled int

	ListTogglesFn     func(ctx context.Context, req togglr.ListTogglesReq) ([]togglr.Toggle, error)
	ListTogglesCalled int

	DeleteToggleFn     func(ctx context.Context, id uid.UID) error
	DeleteToggleCalled int

	Error error
}

func NewToggleService(err error) *ToggleService {
	return &ToggleService{Error: err}
}

func (m *ToggleService) CreateToggle(ctx context.Context, toggle togglr.Toggle) (uid.UID, error) {
	m.CreateToggleCalled++
	if m.CreateToggleFn != nil {
		return m.CreateToggleFn(ctx, toggle)
	}

	if toggle.ID.IsNull() {
		return uid.New(), m.Error
	}

	return toggle.ID, m.Error
}

func (m *ToggleService) UpdateToggle(ctx context.Context, req togglr.UpdateToggleReq) error {
	m.UpdateToggleCalled++
	if m.UpdateToggleFn != nil {
		return m.UpdateToggleFn(ctx, req)
	}

	return m.Error
}

func (m *ToggleService) FetchToggle(ctx context.Context, id uid.UID) (togglr.Toggle, error) {
	m.FetchToggleCalled++
	if m.FetchToggleFn != nil {
		return m.FetchToggleFn(ctx, id)
	}

	return togglr.Toggle{}, m.Error
}

func (m *ToggleService) ListToggles(ctx context.Context, req togglr.ListTogglesReq) ([]togglr.Toggle, error) {
	m.ListTogglesCalled++
	if m.ListTogglesFn != nil {
		return m.ListTogglesFn(ctx, req)
	}

	return make([]togglr.Toggle, 0), m.Error
}

func (m *ToggleService) DeleteToggle(ctx context.Context, id uid.UID) error {
	m.DeleteToggleCalled++
	if m.DeleteToggleFn != nil {
		return m.DeleteToggleFn(ctx, id)
	}

	return m.Error
}
