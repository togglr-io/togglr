package mock

import (
	"context"

	"github.com/togglr-io/togglr"
	"github.com/togglr-io/togglr/uid"
)

type AccountService struct {
	CreateAccountFn     func(ctx context.Context, account togglr.Account) (uid.UID, error)
	CreateAccountCalled int

	UpdateAccountFn     func(ctx context.Context, req togglr.UpdateAccountReq) error
	UpdateAccountCalled int

	FetchAccountFn     func(ctx context.Context, id uid.UID) (togglr.Account, error)
	FetchAccountCalled int

	ListAccountsFn     func(ctx context.Context, req togglr.ListAccountsReq) ([]togglr.Account, error)
	ListAccountsCalled int

	DeleteAccountFn     func(ctx context.Context, id uid.UID) error
	DeleteAccountCalled int

	UpdateAccountUsersFn     func(ctx context.Context, accountID uid.UID, req togglr.UpdateAccountUsersReq) error
	UpdateAccountUsersCalled int

	Error error
}

func NewAccountService(err error) *AccountService {
	return &AccountService{Error: err}
}

func (m *AccountService) CreateAccount(ctx context.Context, account togglr.Account) (uid.UID, error) {
	m.CreateAccountCalled++
	if m.CreateAccountFn != nil {
		return m.CreateAccountFn(ctx, account)
	}

	if account.ID.IsNull() {
		return uid.New(), m.Error
	}

	return account.ID, m.Error
}

func (m *AccountService) UpdateAccount(ctx context.Context, req togglr.UpdateAccountReq) error {
	m.UpdateAccountCalled++
	if m.UpdateAccountFn != nil {
		return m.UpdateAccountFn(ctx, req)
	}

	return m.Error
}

func (m *AccountService) FetchAccount(ctx context.Context, id uid.UID) (togglr.Account, error) {
	m.FetchAccountCalled++
	if m.FetchAccountFn != nil {
		return m.FetchAccountFn(ctx, id)
	}

	return togglr.Account{}, m.Error
}

func (m *AccountService) ListAccounts(ctx context.Context, req togglr.ListAccountsReq) ([]togglr.Account, error) {
	m.ListAccountsCalled++
	if m.ListAccountsFn != nil {
		return m.ListAccountsFn(ctx, req)
	}

	return make([]togglr.Account, 0), m.Error
}

func (m *AccountService) DeleteAccount(ctx context.Context, id uid.UID) error {
	m.DeleteAccountCalled++
	if m.DeleteAccountFn != nil {
		return m.DeleteAccountFn(ctx, id)
	}

	return m.Error
}

func (m *AccountService) UpdateAccountUsers(ctx context.Context, accountID uid.UID, req togglr.UpdateAccountUsersReq) error {
	m.UpdateAccountUsersCalled++
	if m.UpdateAccountUsersFn != nil {
		return m.UpdateAccountUsersFn(ctx, accountID, req)
	}

	return m.Error
}
