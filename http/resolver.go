package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/togglr-io/togglr"
	"github.com/togglr-io/togglr/rules"
	"github.com/togglr-io/togglr/uid"
	"go.uber.org/zap"
)

// HandleResolvePost handles POST requests to the /resolve endpoint
func HandleResolvePost(log *zap.Logger, resolver togglr.Resolver) http.HandlerFunc {
	log = log.With(zap.String("handler", "handleResolverPost"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountID := chi.URLParam(r, "accountID")
		log = log.With(zap.String("accountID", accountID))
		log.Debug("resolving toggles")
		defer log.Sync()

		accountUID, err := uid.FromString(accountID)
		if err != nil {
			log.Error("failed to parse account ID", zap.Error(err))
			badRequest(w, "account ID was badly formed")
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("failed to read request", zap.Error(err))
			serverError(w, "could not read request")
			return
		}

		var rawMetadata map[string]interface{}
		if err := json.Unmarshal(body, &rawMetadata); err != nil {
			log.Error("failed to unmarshal metadata from request", zap.Error(err))
			badRequest(w, "could not unmarshal metadata")
			return
		}

		resolved, err := resolver.Resolve(r.Context(), accountUID, rules.MetaFromRaw(rawMetadata))
		if err != nil {
			log.Error("failed to resolve toggles", zap.Error(err))
			serverError(w, "could not resolve toggles")
			return
		}

		data, err := json.Marshal(resolved)
		if err != nil {
			log.Error("failed to marshal response", zap.Error(err))
			serverError(w, "could not resolve toggles")
			return
		}

		ok(w, data)
	})
}
