package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/togglr-io/togglr"
	"github.com/togglr-io/togglr/uid"
)

// CreateUser creates a new User in postgres
func (c Client) CreateUser(ctx context.Context, user togglr.User) (uid.UID, error) {
	// if no ID is provided, generate one
	if user.ID.IsNull() {
		user.ID = uid.New()
	}

	query := c.db.Insert("users").Rows(user)
	if _, err := query.Executor().ExecContext(ctx); err != nil {
		return user.ID, err
	}

	return user.ID, nil
}

func (c Client) UpdateUser(ctx context.Context, req togglr.UpdateUserReq) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	query := tx.Update("users").Set(updateReqToRecord(req)).Where(goqu.Ex{"id": req.ID})
	if _, err := query.Executor().ExecContext(ctx); err != nil {
		if rbErr := tx.Rollback(); err != nil {
			return fmt.Errorf("rollback failed with err: %s %w", rbErr, err)
		}

		return fmt.Errorf("failed with rollback: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (c Client) FetchUser(ctx context.Context, id uid.UID) (togglr.User, error) {
	var user togglr.User
	query := c.db.From("users").Where(goqu.Ex{"id": id})
	if _, err := query.ScanStructContext(ctx, &user); err != nil {
		return user, err
	}

	return user, nil
}

func (c Client) FetchUserByIdentity(ctx context.Context, identity string, identityType togglr.IdentityType) (togglr.User, error) {
	var user togglr.User
	query := c.db.From("users").Where(goqu.Ex{"identity": identity, "identity_type": identityType})
	if _, err := query.ScanStructContext(ctx, &user); err != nil {
		return user, err
	}

	return user, nil
}

// ListUsers queries a slice of Users from postgres
func (c Client) ListUsers(ctx context.Context, req togglr.ListUsersReq) ([]togglr.User, error) {
	// default to instantiated value so that we return an empty slice instead of null when there's no results
	if req.AccountID.IsNull() {
		return nil, errors.New("attempted to list users without specifiying an account")
	}

	users := []togglr.User{}
	query := c.db.
		Select("u.*").
		From(goqu.T("users").As("u")).
		LeftJoin(
			goqu.T("account_users").As("au"),
			goqu.On(
				goqu.Ex{
					"au.user_id": goqu.I("u.id"),
				},
			),
		)

	if err := query.Where(goqu.Ex{"au.account_id": req.AccountID}).ScanStructsContext(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// DeleteUser deletes a User from postgres
func (c Client) DeleteUser(ctx context.Context, id uid.UID) error {
	del := c.db.Delete("users").Where(goqu.Ex{"id": id}).Executor()
	if _, err := del.Exec(); err != nil {
		return err
	}

	return nil
}
