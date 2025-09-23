package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/ws", s.websocketHandler)

	mux.HandleFunc("POST /user", s.createUser)
	mux.Handle("GET /user/me", s.requireAuthenticatedUser(http.HandlerFunc(s.getUserData)))
	mux.HandleFunc("POST /token", s.createAuthenticationToken)
	mux.Handle("POST /game/new", s.requireAuthenticatedUser(http.HandlerFunc(s.createNewGame)))

	mws := s.CreateMWStack(s.corsMW, s.logAccessMW, s.recoverPanicMW, s.authenticate)
	// Wrap the mux with CORS middleware
	return mws(mux)
}
