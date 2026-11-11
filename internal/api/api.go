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

	mux.HandleFunc("GET /health")

	mux.HandleFunc("POST /users")
}
