package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.healthHandler)
	mux.Handle("/websocket", s.requireAuthenticatedUser(http.HandlerFunc(s.websocketHandler)))

	mux.HandleFunc("POST /users", s.createUser)
	mux.HandleFunc("POST /token", s.createAuthenticationToken)

	mws := s.CreateMWStack(s.corsMW, s.logAccessMW, s.recoverPanicMW, s.authenticate)
	// Wrap the mux with CORS middleware
	return mws(mux)
}
