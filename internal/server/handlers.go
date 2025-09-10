package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"fmt"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/ninox14/gore-codenames/internal/database/sqlc"
	"github.com/ninox14/gore-codenames/internal/request"
	"github.com/ninox14/gore-codenames/internal/response"
	"github.com/ninox14/gore-codenames/internal/validator"
	"github.com/pascaldekloe/jwt"
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

func (s *Server) createAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}
	var v validator.Validator

	err := request.DecodeJSON(w, r, &input)
	s.logger.Log(r.Context(), slog.LevelDebug, "Parsed", "input", input)
	if err != nil {
		s.badRequest(w, r, err)
		return
	}

	v.CheckField(input.Name != "", "Name", "Name is required")
	v.CheckField(input.ID != uuid.Nil, "Id", "ID is missing")

	if v.HasErrors() {
		s.failedValidation(w, r, v)
		return
	}

	user, err := s.db.Queries.GetUserByID(r.Context(), input.ID)
	if err != nil {
		s.serverError(w, r, err)
		return
	}

	var claims jwt.Claims
	claims.Subject = user.ID.String()

	expiry := time.Now().Add(24 * time.Hour)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)

	claims.Issuer = s.config.baseURL
	claims.Audiences = []string{s.config.baseURL}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(s.config.jwt.secretKey))
	if err != nil {
		s.serverError(w, r, err)
		return
	}

	data := map[string]string{
		"AuthenticationToken":       string(jwtBytes),
		"AuthenticationTokenExpiry": expiry.Format(time.RFC3339),
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		s.serverError(w, r, err)
	}
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
