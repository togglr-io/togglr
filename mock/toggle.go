package mock

import (
	"context"

	"github.com/eriktate/toggle"
	"github.com/eriktate/toggle/uid"
)

type ToggleService struct {
	CreateToggleFn     func(ctx context.Context, toggle toggle.Toggle) (uid.UID, error)
	CreateToggleCalled int

	FetchToggleFn     func(ctx context.Context, id uid.UID) (toggle.Toggle, error)
	FetchToggleCalled int

	ListTogglesFn     func(ctx context.Context, req toggle.ListTogglesReq) ([]toggle.Toggle, error)
	ListTogglesCalled int

	Error error
}

func NewToggleService(err error) *ToggleService {
	return &ToggleService{Error: err}
}

func (m *ToggleService) CreateToggle(ctx context.Context, toggle toggle.Toggle) (uid.UID, error) {
	if m.CreateToggleFn != nil {
		return m.CreateToggleFn(ctx, toggle)
	}

	return uid.New(), m.Error
}

func (m *ToggleService) FetchToggle(ctx context.Context, id uid.UID) (toggle.Toggle, error) {
	if m.FetchToggleFn != nil {
		return m.FetchToggleFn(ctx, id)
	}

	return toggle.Toggle{}, m.Error
}

func (m *ToggleService) ListToggles(ctx context.Context, req toggle.ListTogglesReq) ([]toggle.Toggle, error) {
	if m.ListTogglesFn != nil {
		return m.ListTogglesFn(ctx, req)
	}

	return make([]toggle.Toggle, 0), m.Error
}
