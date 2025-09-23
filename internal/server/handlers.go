package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
	"github.com/ninox14/gore-codenames/internal/database/lib"
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
	var v validator.Validator

	err := request.DecodeJSON(w, r, &input)

	if err != nil {
		s.badRequest(w, r, err)
		return
	}

	v.CheckField(input.Name != "", "Name", "Name is required")

	if v.HasErrors() {
		s.failedValidation(w, r, v)
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

	expiry := time.Now().Add(7 * 24 * time.Hour)
	claims.Issued = jwt.NewNumericTime(time.Now().Round(time.Second))
	claims.NotBefore = jwt.NewNumericTime(time.Now().Round(time.Second))
	claims.Expires = jwt.NewNumericTime(expiry.Round(time.Second))

	claims.Issuer = s.config.baseURL
	claims.Audiences = []string{s.config.baseURL}
	s.logger.Log(r.Context(), slog.LevelDebug, "Current secret", "secret", s.config.jwt.secretKey, "length", len([]byte(s.config.jwt.secretKey))*8)

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

func (s *Server) getUserData(w http.ResponseWriter, r *http.Request) {
	user, ok := contextGetAuthenticatedUser(r)
	if !ok {
		s.serverError(w, r, fmt.Errorf("failed to retrieve user data from request context"))
		return
	}
	response.JSON(w, http.StatusOK, user)
}

func (s *Server) createNewGame(w http.ResponseWriter, r *http.Request) {
	user, ok := contextGetAuthenticatedUser(r)
	if !ok {
		s.serverError(w, r, errors.New("failed to retrieve user data from request context"))
		return
	}
	initGameState, err := GetInitialGameState(r.Context(), &user, s.db, s.logger)
	if err != err {
		s.serverError(w, r, err)
		return
	}
	gameId := uuid.New()

	_, err = s.db.Queries.CreateGame(r.Context(), sqlc.CreateGameParams{
		ID:         gameId,
		HostID:     user.ID,
		WordPackID: initGameState.WordPackID,
		GameState:  initGameState,
	})

	if err != nil {
		s.serverError(w, r, err)
		lib.QuietDeleteGame(r.Context(), s.db.Queries, gameId)
		return
	}

	redisKey := fmt.Sprintf("game:%s", gameId)

	// TODO: change redis game expiry
	err = s.rdb.JSONSet(r.Context(), redisKey, "$", initGameState).Err()
	if err != nil {
		s.serverError(w, r, err)
		return
	}

	resp := struct {
		GameID uuid.UUID `json:"game_id"`
	}{
		GameID: gameId,
	}

	response.JSON(w, http.StatusOK, resp)
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := contextGetAuthenticatedUser(r)
	if !ok {
		s.logger.Error("Could not get authed user")
		s.invalidAuthenticationToken(w, r)
		return
	}

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// FIXME: add origin check on deploy
		InsecureSkipVerify: true,
	})

	if err != nil {
		s.serverError(w, r, err)
		return
	}

	ctx := context.Background()
	defer c.Close(websocket.StatusGoingAway, "Normal closure")

	go websocketPingLoop(ctx, c, user.ID, *s.logger)

	for {
		var msg Message
		err := wsjson.Read(ctx, c, &msg)

		switch websocket.CloseStatus(err) {
		case websocket.StatusNormalClosure, websocket.StatusGoingAway:
			return
		}
		if err != nil {
			s.logger.Error("JSON unmarshal error", "error", err)
			break
		}

		s.logger.Debug("Incoming message", "message", msg)
		processWSMessage(ctx, &msg, c, user, s.gh)
	}
}
