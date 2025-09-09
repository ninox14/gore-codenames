package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/websocket", s.websocketHandler)

	mux.HandleFunc("POST /users", s.createUser)

	mws := s.CreateMWStack(s.corsMW, s.logAccessMW, s.recoverPanicMW)
	// Wrap the mux with CORS middleware
	return mws(mux)
}
