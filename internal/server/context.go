package server

import (
	"context"
	"net/http"

	"github.com/ninox14/gore-codenames/internal/database/sqlc"
)

type contextKey string

const (
	authenticatedUserContextKey = contextKey("authenticatedUser")
)

func contextSetAuthenticatedUser(r *http.Request, user sqlc.User) *http.Request {
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
	return r.WithContext(ctx)
}

func contextGetAuthenticatedUser(r *http.Request) (sqlc.User, bool) {
	user, ok := r.Context().Value(authenticatedUserContextKey).(sqlc.User)
	return user, ok
}
