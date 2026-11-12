package api

import (
	"net/http"

	"github.com/jozuacydsagradojuariojcsj/go_api_practice/internal/store"
)

type API struct {
	store *store.Store
}

func New(s *store.Store) *API {
	return &API{store: s}
}

func (a *API) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", a.health)

	mux.HandleFunc("POST /users", a.createUser)
	mux.HandleFunc("GET /users", a.listUsers)
	mux.HandleFunc("GET /users/{id}", a.getUser)
	mux.HandleFunc("PUT /users/{id}", a.updateUser)
	mux.HandleFunc("DELETE /users/{id}", a.deleteUser)

	mux.HandleFunc("POST /users/{id}/deposit", a.deposit)
	mux.HandleFunc("POST /transfers", a.transfer)
	mux.HandleFunc("GET /users/{id}/transactions", a.listUserTransaction)

	return mux
}

func (a *API) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
