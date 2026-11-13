package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

const tokenTTL = 24 * time.Hour

var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("not-a-real-password"), bcrypt.DefaultCost)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	user, err := a.store.GetUserByEmail(r.Context(), req.Email)

	hash := dummyHash
	hasPassword := false
	switch {
	case err == nil:
		if user.PasswordHash != nil {
			hash = []byte(*user.PasswordHash)
			hasPassword = true
		}
	case errors.Is(err, pgx.ErrNoRows):
	default:
		writeError(w, http.StatusInternalServerError, "could not log in")
		return
	}

	cmpErr := bcrypt.CompareHashAndPassword(hash, []byte(req.Password))
	if !hasPassword || cmpErr != nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	token, expires, err := a.signToken(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not log ing")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"token":      token,
		"expires_at": expires,
	})
}

func (a *API) signToken(userID uuid.UUID) (string, time.Time, error) {
	now := time.Now()
	expires := now.Add(tokenTTL)

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expires),
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(a.jwtSecret)
	return signed, expires, err
}

type ctxKey string

const userIDKey ctxKey = "userID"

func (a *API) requireAuth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		tokenStr, found := strings.CutPrefix(header, "Bearer ")
		if !found {
			writeError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{},
			func(t *jwt.Token) (any, error) { return a.jwtSecret, nil },
			jwt.WithValidMethods([]string{"HS256"}),
		)
		if err != nil || !token.Valid {
			log.Printf("invalid token %v", err)
			writeError(w, http.StatusUnauthorized, "invalid token or expired")
			return
		}
		claims := token.Claims.(*jwt.RegisteredClaims)
		id, err := uuid.Parse(claims.Subject)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token subject")
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func currentUserID(r *http.Request) uuid.UUID {
	id, _ := r.Context().Value(userIDKey).(uuid.UUID)
	return id
}
