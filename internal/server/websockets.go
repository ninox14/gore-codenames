package server

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
	"github.com/ninox14/gore-codenames/internal/database"
)

type MessageType string

type GameID uuid.UUID

const (
	MsgJoinGame     MessageType = "join_game"
	MsgLeaveGame    MessageType = "leave_game"
	MsgGameState    MessageType = "game_state"
	MsgPlayerJoined MessageType = "player_joined"
	MsgPlayerLeft   MessageType = "player_left"
)

type Message struct {
	Type   MessageType `json:"type"`
	Data   any         `json:"data"`
	GameID GameID      `json:"game_id,omitempty"`
}

type Player struct {
	ID       uuid.UUID       `json:"id"`
	Name     string          `json:"name"`
	Conn     *websocket.Conn `json:"-"`
	GameID   GameID          `json:"-"`
	LastSeen time.Time       `json:"-"`
}

type Game struct {
	ID      GameID
	Players map[uuid.UUID]*Player
	mu      sync.RWMutex
	hub     *GameHub
}

func NewGame(id GameID, hub *GameHub) *Game {
	return &Game{
		ID:      id,
		Players: make(map[uuid.UUID]*Player),
		hub:     hub,
	}
}

func (g *Game) broadcast(ctx context.Context, msg Message) {
	for _, player := range g.Players {
		if player.Conn != nil {
			if err := wsjson.Write(ctx, player.Conn, msg); err != nil {
				g.hub.logger.Error("Error broadcasting to player", "player", player.ID, "error", err)
				// Remove player if connection is broken
				// go g.RemovePlayer(player.ID)
			}
		}
	}
}

func (g *Game) AddPlayer(player *Player) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Players[player.ID] = player
	player.GameID = g.ID
	// player.LastSeen = time.Now()

	// Broadcast updated game state to all players in lobby
	// g.broadcastGameState()
}

func (g *Game) RemovePlayer(playerID uuid.UUID) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if player, exists := g.Players[playerID]; exists {
		delete(g.Players, playerID)

		// Notify other players that this player left
		// g.broadcast(Message{
		// 	Type: MsgPlayerLeft,
		// 	Data: map[string]string{"player_id": playerID},
		// })

		// Close the connection
		if player.Conn != nil {
			player.Conn.Close(websocket.StatusNormalClosure, "Player Left the game")
		}

		// If lobby is empty, remove it
		if len(g.Players) == 0 {
			g.hub.RemoveLobby(g.ID)
		}
	}
}

type GameHub struct {
	games  map[GameID]*Game
	mu     sync.RWMutex
	logger *slog.Logger
	db     *database.DB
}

func NewGameHub(logger *slog.Logger, db *database.DB) *GameHub {
	return &GameHub{games: make(map[GameID]*Game), logger: logger, db: db}
}

func (h *GameHub) RemoveLobby(gameID GameID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if game, exists := h.games[gameID]; exists && len(game.Players) == 0 {
		delete(h.games, gameID)
		h.logger.Debug("Removed empty game", "gameId", gameID)
	}
}

func websocketPingLoop(ctx context.Context, c *websocket.Conn, userId uuid.UUID, logger slog.Logger) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

			pingCtx, pingCancel := context.WithTimeout(ctx, 10*time.Second)
			err := c.Ping(pingCtx)
			pingCancel()

			if err != nil {
				logger.Error("Failed ping for", "userId", userId.String())
				c.Close(websocket.StatusAbnormalClosure, "Failed ping response")
				return
			}
		}
	}
}

func processWSMessage(ctx context.Context, msg *Message, c *websocket.Conn) {
	switch msg.Type {
	case MsgJoinGame:
		wsjson.Write(ctx, c, msg)
	default:
		wsjson.Write(ctx, c, msg)
	}
}
