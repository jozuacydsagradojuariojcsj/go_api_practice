package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jozuacydsagradojuariojcsj/go_api_practice/internal/db"
)

type depositRequest struct {
	Amount int64 `json:"amount"`
}

type transferRequest struct {
	From   uuid.UUID `json:"from"`
	To     uuid.UUID `json:"to"`
	Amount int64     `json:"amount"`
}

func (a *API) deposit(w http.ResponseWriter, r *http.Request) {
	id := currentUserID(r)

	var req depositRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
	}

	if err := a.store.Deposit(r.Context(), id, req.Amount); err != nil {
		writeStoreError(w, err)
		return
	}

	user, err := a.store.GetUser(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toUserResponse(user))
}

func (a *API) transfer(w http.ResponseWriter, r *http.Request) {
	from := currentUserID(r)

	var req transferRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON Body")
		return
	}

	if err := a.store.Transfer(r.Context(), from, req.To, req.Amount); err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "transfer complete"})
}

func (a *API) listUserTransaction(w http.ResponseWriter, r *http.Request) {
	id := currentUserID(r)

	txs, err := a.store.ListUserTransaction(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if txs == nil {
		txs = []db.Transaction{}
	}
	writeJSON(w, http.StatusOK, txs)
}
