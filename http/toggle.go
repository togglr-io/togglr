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

// HandleTogglePOST handles POST requests to the /toggle endpoint
func HandleTogglePOST(logger *zap.Logger, ts togglr.ToggleService) http.HandlerFunc {
	log := logger.With(zap.String("handler", "HandleTogglePOST"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("creating toggle")

		var id togglr.ID
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("failed to read request", zap.Error(err))
			serverError(w, "could not read request")
			return
		}

		if err := json.Unmarshal(body, &id); err != nil {
			log.Error("failed to unmarshal ID from request", zap.Error(err))
			badRequest(w, "could not unmarshal toggle")
			return
		}

		// whether or not the Toggle has an ID determines if we're creating a new one or
		// updating an existing one
		if id.ID.IsNull() {
			var toggle togglr.Toggle
			if err := json.Unmarshal(body, &toggle); err != nil {
				log.Error("failed to unmarshal toggle", zap.Error(err))
				badRequest(w, "could not unmarshal toggle")
				return
			}

			id.ID, err = ts.CreateToggle(r.Context(), toggle)
			if err != nil {
				log.Error("failed to create toggle", zap.Error(err))
				serverError(w, "could not save toggle")
				return
			}
		} else {
			var updateReq togglr.UpdateToggleReq
			if err := json.Unmarshal(body, &updateReq); err != nil {
				log.Error("failed to unmarshal update req", zap.Error(err))
				badRequest(w, "could not unmarshal toggle")
				return
			}

			if err := ts.UpdateToggle(r.Context(), updateReq); err != nil {
				log.Error("failed to update toggle", zap.Error(err))
				serverError(w, "could not save toggle")
				return
			}
		}

		data, err := json.Marshal(id)
		if err != nil {
			log.Error("failed to marshal response", zap.Error(err))
			serverError(w, "could not create toggle")
			return
		}

		ok(w, data)
	})
}

// HandleToggleGET handles GET requests to the /toggle endpoint
func HandleToggleGET(logger *zap.Logger, ts togglr.ToggleService) http.HandlerFunc {
	log := logger.With(zap.String("handler", "HandleToggleGET"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("listing toggles")

		req := togglr.ListTogglesReq{}

		toggles, err := ts.ListToggles(r.Context(), req)
		if err != nil {
			log.Error("failed to list toggles", zap.Error(err))
			serverError(w, "could not list toggles")
			return
		}

		data, err := json.Marshal(toggles)
		if err != nil {
			log.Error("failed to marshal toggles", zap.Error(err))
			serverError(w, "could not list toggles")
			return
		}

		ok(w, data)
	})
}

// HandleToggleGetIdGET handles GET requests to the /toggle/{id} endpoint
func HandleToggleIdGET(logger *zap.Logger, ts togglr.ToggleService) http.HandlerFunc {
	log := logger.With(zap.String("handler", "HandleToggleIdGET"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		log = log.With(zap.String("toggleID", id))
		log.Debug("fetching toggle")

		uid, err := uid.FromString(id)
		if err != nil {
			log.Error("failed to parse toggle ID", zap.Error(err))
			badRequest(w, "toggle ID was badly formed")
			return
		}

		tog, err := ts.FetchToggle(r.Context(), uid)
		if err != nil {
			log.Error("failed to fetch toggle", zap.Error(err))
			serverError(w, "could not fetch toggle")
			return
		}

		data, err := json.Marshal(tog)
		if err != nil {
			log.Error("failed to marshal toggle", zap.Error(err))
			serverError(w, "could not fetch toggle")
			return
		}

		ok(w, data)
	})
}

// HandleToggleDELETE handles DELETE requests to the /toggle/{id} endpoint
func HandleToggleDELETE(logger *zap.Logger, ts togglr.ToggleService) http.HandlerFunc {
	log := logger.With(zap.String("handler", "HandleToggleDELETE"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		log = log.With(zap.String("toggleID", id))
		log.Debug("deleting toggle")

		uid, err := uid.FromString(id)
		if err != nil {
			log.Error("failed to parse toggle ID", zap.Error(err))
			badRequest(w, "toggle ID was badly formed")
			return
		}

		if err := ts.DeleteToggle(r.Context(), uid); err != nil {
			log.Error("failed to delete toggle", zap.Error(err))
			serverError(w, "could not delete toggle")
			return
		}

		noContent(w)
	})
}
