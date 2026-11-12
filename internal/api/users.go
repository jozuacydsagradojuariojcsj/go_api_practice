package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jozuacydsagradojuariojcsj/go_api_practice/internal/db"
)

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (a *API) createUser(w http.ResponseWriter, r *http.Request){
	var req createUserRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON Body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	if req.Name == "" || req.Email == ""{
		writeError(w, http.StatusBadRequest, "name and email are required")
		return
	}

	user,err := a.store.CreateUser(r.Context(), db.CreateUserParams{
		Name: req.Name,
		Email: req.Email,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505"{
			writeError(w, http.StatusConflict,"email already exists")
			return
		}
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (a *API) getUser(w http.ResponseWriter, r *http.Request){
	id, ok := pathUUID(w, r, "id")
	if !ok {
		return
	}

	user, err := a.store.GetUser(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (a *API) listUsers(w http.ResponseWriter, r *http.Request){
	users, err := a.store.ListUsers(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}

	if users == nil {
		users = []db.User{}
	}
	writeJSON(w, http.StatusOK, users)
}