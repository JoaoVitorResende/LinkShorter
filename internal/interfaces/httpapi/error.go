package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"shortapi/internal/exception"
)

var errInvalidBody = errors.New("invalid body")

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, exception.ErrNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
	case errors.Is(err, exception.ErrEmptyURL),
		errors.Is(err, exception.ErrEmptyCode),
		errors.Is(err, exception.ErrInvalidURL),
		errors.Is(err, errInvalidBody):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	case errors.Is(err, exception.ErrCodeTaken):
		writeJSON(w, http.StatusConflict, errorResponse{Error: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "erro desconhecido"})
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}