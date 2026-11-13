package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jozuacydsagradojuariojcsj/go_api_practice/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userUpdateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func toUserResponse(u db.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Balance:   u.Balance,
		CreatedAt: u.CreatedAt,
	}
}

func (a *API) createUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON Body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	if req.Name == "" || req.Email == "" {
		writeError(w, http.StatusBadRequest, "name and email are required")
		return
	}
	if len(req.Password) < 8 || len(req.Password) > 72 {
		writeError(w, http.StatusBadRequest, "password must be 8 to 72 characters")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	hashStr := string(hash)

	user, err := a.store.CreateUser(r.Context(), db.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: &hashStr,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeError(w, http.StatusConflict, "email already exists")
			return
		}
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toUserResponse(user))
}

func (a *API) getUser(w http.ResponseWriter, r *http.Request) {
	id := currentUserID(r)

	user, err := a.store.GetUser(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (a *API) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.store.ListUsers(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}

	if users == nil {
		users = []db.User{}
	}
	out := make([]userResponse, 0, len(users))
	for _, u := range users {
		out = append(out, toUserResponse(u))
	}
	writeJSON(w, http.StatusOK, out)
}

func (a *API) updateUser(w http.ResponseWriter, r *http.Request) {
	id := currentUserID(r)

	var req userUpdateRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	if req.Name == "" || req.Email == "" {
		writeError(w, http.StatusBadRequest, "username or email required")
		return
	}

	var hashPtr *string
	if req.Password != "" {
		if len(req.Password) < 8 || len(req.Password) > 72 {
			writeError(w, http.StatusBadRequest, "password must be between 8 to 72 characters")
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not update user")
		}

		h := string(hash)
		hashPtr = &h
	}

	user, err := a.store.UpdateUser(r.Context(), db.UpdateUserParams{
		ID:           id,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashPtr,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeError(w, http.StatusConflict, "email already exists")
			return
		}
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toUserResponse(user))
}

func (a *API) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := currentUserID(r)
	n, err := a.store.DeleteUser(r.Context(), id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			writeError(w, http.StatusConflict, "user has transaction history cannot be deleted")
			return
		}
		writeStoreError(w, err)
		return
	}

	if n == 0 {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
