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

func handleTogglePost(log *zap.Logger, ts toggle.ToggleService) http.HandlerFunc {
	log = log.With(zap.String("handler", "handleTogglePost"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("creating toggle")
		defer log.Sync()

		var tog toggle.Toggle
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("failed to read request", zap.Error(err))
			serverError(w, "could not read request")
			return
		}

		if err := json.Unmarshal(data, &tog); err != nil {
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
		data, err = json.Marshal(res)
		if err != nil {
			log.Error("failed to marshal response", zap.Error(err))
			serverError(w, "could not create toggle")
			return
		}

		ok(w, data)
	})
}

func handleToggleList(log *zap.Logger, ts toggle.ToggleService) http.HandlerFunc {
	log = log.With(zap.String("handler", "handleToggleList"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer log.Sync()

		noContent(w)
	})
}

func handleToggleDetail(log *zap.Logger, ts toggle.ToggleService) http.HandlerFunc {
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

func handleToggleDelete(log *zap.Logger, ts toggle.ToggleService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		w.Write(nil)
	})
}
