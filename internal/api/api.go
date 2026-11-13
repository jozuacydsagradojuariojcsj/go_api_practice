package api

import (
	"net/http"

	"github.com/jozuacydsagradojuariojcsj/go_api_practice/internal/store"
)

type API struct {
	store     *store.Store
	jwtSecret []byte
}

func New(s *store.Store, jwtSecret []byte) *API {
	return &API{store: s, jwtSecret: jwtSecret}
}

func (a *API) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", a.health)
	mux.HandleFunc("POST /users", a.createUser)
	mux.HandleFunc("POST /auth/login", a.login)
	mux.HandleFunc("GET /users", a.listUsers)

	mux.Handle("GET /user/me", a.requireAuth(a.getUser))
	mux.Handle("PUT /users", a.requireAuth(a.updateUser))
	mux.Handle("DELETE /users", a.requireAuth(a.deleteUser))

	mux.Handle("POST /deposit", a.requireAuth(a.deposit))
	mux.Handle("POST /transfers", a.requireAuth(a.transfer))
	mux.Handle("GET /transactions", a.requireAuth(a.listUserTransaction))

	return mux
}

func (a *API) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
