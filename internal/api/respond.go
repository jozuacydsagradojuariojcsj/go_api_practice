package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jozuacydsagradojuariojcsj/go_api_practice/internal/store"
)

func decodeJSON(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1 << 20)
	return json.NewDecoder(r.Body).Decode(v)
}

func writeJSON(w http.ResponseWriter, status int, v any){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string){
	writeJSON(w, status, map[string]string{"error":msg})
}

func pathUUID(w http.ResponseWriter, r *http.Request, name string)(uuid.UUID, bool){
	id, err := uuid.Parse(r.PathValue(name))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid" + name)
		return uuid.Nil, false
	}
	return id, true
}

func writeStoreError(w http.ResponseWriter, err error){
	switch{
	case errors.Is(err, store.ErrInsufficientFunds):
		writeError(w, http.StatusUnprocessableEntity, "insufficient funds")
	case errors.Is(err, store.ErrInvalidTransfer):
		writeError(w, http.StatusBadRequest, "invalid request: amount must be > 0 and users must differ")
	case errors.Is(err, store.ErrUserNotFound):
		writeError(w, http.StatusNotFound, "User not Found")
	default:
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}