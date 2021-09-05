package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/togglr-io/togglr"
	"github.com/togglr-io/togglr/uid"
	"go.uber.org/zap"
)

// HandleAccountPOST handles POST requests to the /account endpoint
func HandleAccountPOST(log *zap.Logger, as togglr.AccountService) http.HandlerFunc {
	log = log.With(zap.String("handler", "HandleAccountPOST"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("creating account")
		defer log.Sync()

		var id togglr.ID
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("failed to read request", zap.Error(err))
			serverError(w, "could not read request")
			return
		}

		if err := json.Unmarshal(body, &id); err != nil {
			log.Error("failed to unmarshal ID from request", zap.Error(err))
			badRequest(w, "could not unmarshal account")
			return
		}

		// whether or not the Account has an ID determines if we're creating a new one or
		// updating an existing one
		if id.ID.IsNull() {
			var account togglr.Account
			if err := json.Unmarshal(body, &account); err != nil {
				log.Error("failed to unmarshal account", zap.Error(err))
				badRequest(w, "could not unmarshal account")
				return
			}

			id.ID, err = as.CreateAccount(r.Context(), account)
			if err != nil {
				log.Error("failed to create account", zap.Error(err))
				serverError(w, "could not save account")
				return
			}
		} else {
			var updateReq togglr.UpdateAccountReq
			if err := json.Unmarshal(body, &updateReq); err != nil {
				log.Error("failed to unmarshal update req", zap.Error(err))
				badRequest(w, "could not unmarshal account")
				return
			}

			if err := as.UpdateAccount(r.Context(), updateReq); err != nil {
				log.Error("failed to update account", zap.Error(err))
				serverError(w, "could not save account")
				return
			}
		}

		data, err := json.Marshal(id)
		if err != nil {
			log.Error("failed to marshal response", zap.Error(err))
			serverError(w, "could not create account")
			return
		}

		ok(w, data)
	})
}

// HandleAccountGET handles GET requests to the /account endpoint
func HandleAccountGET(log *zap.Logger, as togglr.AccountService) http.HandlerFunc {
	log = log.With(zap.String("handler", "HandleAccountGET"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("listing accounts")
		defer log.Sync()

		req := togglr.ListAccountsReq{}

		accounts, err := as.ListAccounts(r.Context(), req)
		if err != nil {
			log.Error("failed to list accounts", zap.Error(err))
			serverError(w, "could not list accounts")
			return
		}

		data, err := json.Marshal(accounts)
		if err != nil {
			log.Error("failed to marshal accounts", zap.Error(err))
			serverError(w, "could not list accounts")
			return
		}

		ok(w, data)
	})
}

// HandleAccountIdGET handles GET requests to the /account/{id} endpoint
func HandleAccountIdGET(log *zap.Logger, as togglr.AccountService) http.HandlerFunc {
	log = log.With(zap.String("handler", "HandleAccountIdGET"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		log = log.With(zap.String("accountID", id))
		log.Debug("fetching account")
		defer log.Sync()

		uid, err := uid.FromString(id)
		if err != nil {
			log.Error("failed to parse account ID", zap.Error(err))
			badRequest(w, "account ID was badly formed")
			return
		}

		account, err := as.FetchAccount(r.Context(), uid)
		if err != nil {
			log.Error("failed to fetch account", zap.Error(err))
			serverError(w, "could not fetch account")
			return
		}

		data, err := json.Marshal(account)
		if err != nil {
			log.Error("failed to marshal account", zap.Error(err))
			serverError(w, "could not fetch account")
			return
		}

		ok(w, data)
	})
}

// HandleAccountUsersGET handles GET requests to the /account/{id} endpoint
func HandleAccountUsersGET(log *zap.Logger, us togglr.UserService) http.HandlerFunc {
	log = log.With(zap.String("handler", "HandleAccountUsersGET"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		log = log.With(zap.String("accountID", id))
		log.Debug("listing account users")
		defer log.Sync()

		uid, err := uid.FromString(id)
		if err != nil {
			log.Error("failed to parse account ID", zap.Error(err))
			badRequest(w, "account ID was badly formed")
			return
		}

		req := togglr.ListUsersReq{
			AccountID: uid,
		}

		users, err := us.ListUsers(r.Context(), req)
		if err != nil {
			log.Error("failed to fetch account users", zap.Error(err))
			serverError(w, "could not fetch account users")
			return
		}

		data, err := json.Marshal(users)
		if err != nil {
			log.Error("failed to marshal account", zap.Error(err))
			serverError(w, "could not fetch account")
			return
		}

		ok(w, data)
	})
}

// HandleAccountUsersPOST handles POST requests to the /account/{id}/user endpoint
func HandleAccountUsersPOST(log *zap.Logger, as togglr.AccountService) http.HandlerFunc {
	log = log.With(zap.String("handler", "HandleAccountUsersPOST"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		log = log.With(zap.String("accountID", id))
		log.Debug("updating account users")
		defer log.Sync()

		accountID, err := uid.FromString(id)
		if err != nil {
			log.Error("failed to parse account ID", zap.Error(err))
			badRequest(w, "account ID was badly formed")
			return
		}

		var req togglr.UpdateAccountUsersReq
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("failed to read request", zap.Error(err))
			serverError(w, "could not read request")
			return
		}

		if err := json.Unmarshal(body, &req); err != nil {
			log.Error("failed to unmarshal user ID from request", zap.Error(err))
			badRequest(w, "could not unmarshal account")
			return
		}

		if err := as.UpdateAccountUsers(r.Context(), accountID, req); err != nil {
			log.Error("failed to add users", zap.Error(err))
			serverError(w, "could not add users to account")
			return
		}

		noContent(w)
	})
}
