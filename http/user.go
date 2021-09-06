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

// HandleUserPOST handles POST requests to the /user endpoint
func HandleUserPOST(logger *zap.Logger, as togglr.UserService) http.HandlerFunc {
	log := logger.With(zap.String("handler", "HandleUserPOST"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("creating user")

		var id togglr.ID
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("failed to read request", zap.Error(err))
			serverError(w, "could not read request")
			return
		}

		if err := json.Unmarshal(body, &id); err != nil {
			log.Error("failed to unmarshal ID from request", zap.Error(err))
			badRequest(w, "could not unmarshal user")
			return
		}

		// whether or not the User has an ID determines if we're creating a new one or
		// updating an existing one
		if id.ID.IsNull() {
			var user togglr.User
			if err := json.Unmarshal(body, &user); err != nil {
				log.Error("failed to unmarshal user", zap.Error(err))
				badRequest(w, "could not unmarshal user")
				return
			}

			id.ID, err = as.CreateUser(r.Context(), user)
			if err != nil {
				log.Error("failed to create user", zap.Error(err))
				serverError(w, "could not save user")
				return
			}
		} else {
			var updateReq togglr.UpdateUserReq
			if err := json.Unmarshal(body, &updateReq); err != nil {
				log.Error("failed to unmarshal update req", zap.Error(err))
				badRequest(w, "could not unmarshal user")
				return
			}

			if err := as.UpdateUser(r.Context(), updateReq); err != nil {
				log.Error("failed to update user", zap.Error(err))
				serverError(w, "could not save user")
				return
			}
		}

		data, err := json.Marshal(id)
		if err != nil {
			log.Error("failed to marshal response", zap.Error(err))
			serverError(w, "could not create user")
			return
		}

		ok(w, data)
	})
}

// HandleUserGET handles GET requests to the /user endpoint
func HandleUserGET(logger *zap.Logger, as togglr.UserService) http.HandlerFunc {
	log := logger.With(zap.String("handler", "HandleUserGET"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("listing users")

		req := togglr.ListUsersReq{}

		users, err := as.ListUsers(r.Context(), req)
		if err != nil {
			log.Error("failed to list users", zap.Error(err))
			serverError(w, "could not list users")
			return
		}

		data, err := json.Marshal(users)
		if err != nil {
			log.Error("failed to marshal users", zap.Error(err))
			serverError(w, "could not list users")
			return
		}

		ok(w, data)
	})
}

// HandleUserIdGET handles GET requests to the /user/{id} endpoint
func HandleUserIdGET(logger *zap.Logger, as togglr.UserService) http.HandlerFunc {
	log := logger.With(zap.String("handler", "HandleUserIdGET"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		log = log.With(zap.String("userID", id))
		log.Debug("fetching user")

		uid, err := uid.FromString(id)
		if err != nil {
			log.Error("failed to parse user ID", zap.Error(err))
			badRequest(w, "user ID was badly formed")
			return
		}

		user, err := as.FetchUser(r.Context(), uid)
		if err != nil {
			log.Error("failed to fetch user", zap.Error(err))
			serverError(w, "could not fetch user")
			return
		}

		data, err := json.Marshal(user)
		if err != nil {
			log.Error("failed to marshal user", zap.Error(err))
			serverError(w, "could not fetch user")
			return
		}

		ok(w, data)
	})
}

// HandleUserDELETE handles DELETE requests to the /user/{id} endpoint
func HandleUserDELETE(logger *zap.Logger, ts togglr.UserService) http.HandlerFunc {
	log := logger.With(zap.String("handler", "HandleUserDELETE"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		log = log.With(zap.String("userID", id))
		log.Debug("deleting user")

		uid, err := uid.FromString(id)
		if err != nil {
			log.Error("failed to parse user ID", zap.Error(err))
			badRequest(w, "user ID was badly formed")
			return
		}

		if err := ts.DeleteUser(r.Context(), uid); err != nil {
			log.Error("failed to delete user", zap.Error(err))
			serverError(w, "could not delete user")
			return
		}

		noContent(w)
	})
}
