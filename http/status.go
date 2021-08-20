package http

import "net/http"

func ok(w http.ResponseWriter, data []byte) {
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func noContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}

func badRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(msg))
}

func forbidden(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(msg))
}

func unauthorized(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(msg))
}

func serverError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}
