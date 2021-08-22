package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/eriktate/toggle"
	"github.com/eriktate/toggle/uid"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// HandleTogglePost handles POST requests to the /toggle endpoint
func HandleTogglePost(log *zap.Logger, ts toggle.ToggleService) http.HandlerFunc {
	log = log.With(zap.String("handler", "handleTogglePost"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("creating toggle")
		defer log.Sync()

		var tog toggle.Toggle
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("failed to read request", zap.Error(err))
			serverError(w, "could not read request")
			return
		}

		if err := json.Unmarshal(body, &tog); err != nil {
			log.Error("failed to unmarshal request", zap.Error(err))
			badRequest(w, "could not unmarshal toggle")
			return
		}

		id, err := ts.CreateToggle(r.Context(), tog)
		if err != nil {
			log.Error("failed to create toggle", zap.Error(err))
			serverError(w, "could not create toggle")
			return
		}

		res := idEnvelope{id}
		data, err := json.Marshal(res)
		if err != nil {
			log.Error("failed to marshal response", zap.Error(err))
			serverError(w, "could not create toggle")
			return
		}

		ok(w, data)
	})
}

// HandleToggleGet handles GET requests to the /toggle endpoint
func HandleToggleGet(log *zap.Logger, ts toggle.ToggleService) http.HandlerFunc {
	log = log.With(zap.String("handler", "handleToggleList"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("listing toggles")
		defer log.Sync()

		req := toggle.ListTogglesReq{}

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

// HandleToggleGetID handles GET requests to the /toggle/{id} endpoint
func HandleToggleGetID(log *zap.Logger, ts toggle.ToggleService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		log.Debug("fetching toggle", zap.String("id", id))
		defer log.Sync()

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
			log.Error("failed to unmarshal toggle", zap.Error(err))
			serverError(w, "could not fetch toggle")
			return
		}

		ok(w, data)
	})
}

// HandleToggleDelete handles DELETE requests to the /toggle/{id} endpoint
func HandleToggleDelete(log *zap.Logger, ts toggle.ToggleService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		log.Debug("deleting toggle", zap.String("id", id))
		defer log.Sync()

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
