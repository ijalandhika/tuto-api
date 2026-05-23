package response

import (
	"encoding/json"
	"net/http"
)

func OK(w http.ResponseWriter, data any) {
	write(w, http.StatusOK, Response{Data: data})
}

func Created(w http.ResponseWriter, data any) {
	write(w, http.StatusCreated, Response{Data: data})
}

func NoContent(w http.ResponseWriter) {
	write(w, http.StatusNoContent, Response{})
}

func BadRequest(w http.ResponseWriter, msg string) {
	write(w, http.StatusBadRequest, Response{Error: msg})
}

func Unauthorized(w http.ResponseWriter, msg string) {
	write(w, http.StatusUnauthorized, Response{Error: msg})
}

func NotFound(w http.ResponseWriter, msg string) {
	write(w, http.StatusNotFound, Response{Error: msg})
}

func InternalServerError(w http.ResponseWriter, msg string) {
	write(w, http.StatusInternalServerError, Response{Error: msg})
}

func write(w http.ResponseWriter, status int, r Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(r)
}
