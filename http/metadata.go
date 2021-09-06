package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/togglr-io/togglr"
	"github.com/togglr-io/togglr/uid"
	"go.uber.org/zap"
)

func HandleMetadataGET(logger *zap.Logger, ms togglr.MetadataService) http.HandlerFunc {
	log := logger.With(zap.String("handler", "handleMetadataGET"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountID := chi.URLParam(r, "accountID")
		log.Debug("listing metdata keys", zap.String("accountID", accountID))

		uid, err := uid.FromString(accountID)
		if err != nil {
			log.Error("failed to parse account ID", zap.Error(err))
			badRequest(w, "account ID was badly formed")
			return
		}

		keys, err := ms.FetchKeys(r.Context(), uid)
		if err != nil {
			log.Error("failed to fetch metdata keys", zap.Error(err))
			serverError(w, "could not fetch metadata keys")
			return
		}

		res := make([]string, len(keys))
		for idx, key := range keys {
			res[idx] = key.Key
		}

		data, err := json.Marshal(res)
		if err != nil {
			log.Error("failed to marshal metadata keys", zap.Error(err))
			serverError(w, "could not fetch metadata keys")
			return
		}

		ok(w, data)
	})
}
