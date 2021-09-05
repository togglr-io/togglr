package mock

import (
	"context"

	"github.com/togglr-io/togglr"
	"github.com/togglr-io/togglr/uid"
)

type UserService struct {
	CreateUserFn     func(ctx context.Context, user togglr.User) (uid.UID, error)
	CreateUserCalled int

	UpdateUserFn     func(ctx context.Context, req togglr.UpdateUserReq) error
	UpdateUserCalled int

	FetchUserFn     func(ctx context.Context, id uid.UID) (togglr.User, error)
	FetchUserCalled int

	ListUsersFn     func(ctx context.Context, req togglr.ListUsersReq) ([]togglr.User, error)
	ListUsersCalled int

	DeleteUserFn     func(ctx context.Context, id uid.UID) error
	DeleteUserCalled int

	Error error
}

func NewUserService(err error) *UserService {
	return &UserService{Error: err}
}

func (m *UserService) CreateUser(ctx context.Context, user togglr.User) (uid.UID, error) {
	m.CreateUserCalled++
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, user)
	}

	if user.ID.IsNull() {
		return uid.New(), m.Error
	}

	return user.ID, m.Error
}

func (m *UserService) UpdateUser(ctx context.Context, req togglr.UpdateUserReq) error {
	m.UpdateUserCalled++
	if m.UpdateUserFn != nil {
		return m.UpdateUserFn(ctx, req)
	}

	return m.Error
}

func (m *UserService) FetchUser(ctx context.Context, id uid.UID) (togglr.User, error) {
	m.FetchUserCalled++
	if m.FetchUserFn != nil {
		return m.FetchUserFn(ctx, id)
	}

	return togglr.User{}, m.Error
}

func (m *UserService) ListUsers(ctx context.Context, req togglr.ListUsersReq) ([]togglr.User, error) {
	m.ListUsersCalled++
	if m.ListUsersFn != nil {
		return m.ListUsersFn(ctx, req)
	}

	return make([]togglr.User, 0), m.Error
}

func (m *UserService) DeleteUser(ctx context.Context, id uid.UID) error {
	m.DeleteUserCalled++
	if m.DeleteUserFn != nil {
		return m.DeleteUserFn(ctx, id)
	}

	return m.Error
}
