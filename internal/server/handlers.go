package server

import (
	"encoding/json"
	"log"
	"net/http"

	"fmt"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/ninox14/gore-codenames/internal/database/sqlc"
	"github.com/ninox14/gore-codenames/internal/request"
	"github.com/ninox14/gore-codenames/internal/response"
	"github.com/ninox14/gore-codenames/internal/validator"
)

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp, err := json.Marshal(s.db.Health(ctx))
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"Name"`
	}
	var v struct {
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)

	if err != nil {
		s.badRequest(w, r, err)
		return
	}

	v.Validator.CheckField(input.Name != "", "Name", "Name is required")

	if v.Validator.HasErrors() {
		s.failedValidation(w, r, v.Validator)
		return
	}

	newId := uuid.New()

	var usrDto = sqlc.CreateUserParams{
		Name: input.Name,
		ID:   newId,
	}

	usr, err := s.db.Queries.CreateUser(r.Context(), usrDto)

	if err != nil {
		s.serverError(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, usr)
}

// TODO: Use socket io
func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	socket, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to open websocket", http.StatusInternalServerError)
		return
	}
	defer socket.Close(websocket.StatusGoingAway, "Server closing websocket")

	ctx := r.Context()
	socketCtx := socket.CloseRead(ctx)

	for {
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		if err := socket.Write(socketCtx, websocket.MessageText, []byte(payload)); err != nil {
			log.Printf("Failed to write to socket: %v", err)
			break
		}
		time.Sleep(2 * time.Second)
	}
}
