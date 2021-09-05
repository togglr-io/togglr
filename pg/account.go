package pg

import (
	"context"
	"fmt"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/togglr-io/togglr"
	"github.com/togglr-io/togglr/uid"
)

// CreateAccount creates a new Account in postgres
func (c Client) CreateAccount(ctx context.Context, account togglr.Account) (uid.UID, error) {
	// if no ID is provided, generate one
	if account.ID.IsNull() {
		account.ID = uid.New()
	}

	query := c.db.Insert("accounts").Rows(account)
	if _, err := query.Executor().ExecContext(ctx); err != nil {
		return account.ID, err
	}

	return account.ID, nil
}

func (c Client) UpdateAccount(ctx context.Context, req togglr.UpdateAccountReq) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	query := tx.Update("accounts").Set(updateReqToRecord(req)).Where(goqu.Ex{"id": req.ID})
	if _, err := query.Executor().ExecContext(ctx); err != nil {
		return c.handleTxErr(tx, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (c Client) FetchAccount(ctx context.Context, id uid.UID) (togglr.Account, error) {
	var account togglr.Account
	query := c.db.From("accounts").Where(goqu.Ex{"id": id})
	if _, err := query.ScanStructContext(ctx, &account); err != nil {
		return account, err
	}

	return account, nil
}

// ListAccounts queries a slice of Accounts from postgres
func (c Client) ListAccounts(ctx context.Context, req togglr.ListAccountsReq) ([]togglr.Account, error) {
	// default to instantiated value so that we return an empty slice instead of null when there's no results
	accounts := []togglr.Account{}
	query := c.db.From("accounts")

	if err := query.ScanStructsContext(ctx, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (c Client) UpdateAccountUsers(ctx context.Context, accountID uid.UID, req togglr.UpdateAccountUsersReq) error {
	// There are 4 possible cases to deal with when updating account users
	// 1. We do nothing (e.g. Add and Remove slices are empty)
	// 2. We _only_ add users
	// 3. We _only_ remove users
	// 4. We add and remove users at the same time

	// first we'll short circuit on the "do nothing" case
	if len(req.Add) == 0 && len(req.Remove) == 0 {
		return nil
	}

	log.Printf("req: %+v", req)
	// we have at least one query to run, so go ahead and create the tx
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// handle the case where there are users to add
	if len(req.Add) > 0 {
		// because goqu tries to be generic, and go doesn't let you spread an array of one
		// type into an array of interface{}, we have to make an array of interface{} and then
		// assign goqu.Records to it (which is dumb, imo)
		rows := make([]interface{}, len(req.Add))
		for idx, userID := range req.Add {
			rows[idx] = goqu.Record{
				"account_id": accountID,
				"user_id":    userID,
			}
		}
		query := tx.Insert("account_users").Rows(rows...).OnConflict(goqu.DoNothing())
		if err != nil {
			log.Printf("failed to generate SQL: %s", err)
			return err
		}
		if _, err := query.Executor().ExecContext(ctx); err != nil {
			return c.handleTxErr(tx, err)
		}
	}

	// handle the case where there are users to remove
	if len(req.Remove) > 0 {
		query := tx.Delete("account_users").Where(goqu.Ex{"account_id": accountID, "user_id": req.Remove})
		if _, err := query.Executor().ExecContext(ctx); err != nil {
			return c.handleTxErr(tx, err)
		}
	}

	// commit all the things!
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}
