package http

import "net/http"

func ok(w http.ResponseWriter, data []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func noContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write(nil)
}

func found(w http.ResponseWriter, url string) {
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
	_, _ = w.Write(nil)
}

func badRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(msg))
}

// nolint
func forbidden(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(msg))
}

// nolint
func unauthorized(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(msg))
}

func serverError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(msg))
}
