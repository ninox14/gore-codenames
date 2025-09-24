package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
	"github.com/ninox14/gore-codenames/internal/database"
	"github.com/ninox14/gore-codenames/internal/database/dto"
	"github.com/ninox14/gore-codenames/internal/database/sqlc"
	"github.com/redis/go-redis/v9"
)

type MessageType string

type RedisPlayersPath string

const (
	SpectatorsPath RedisPlayersPath = "spectators"
	TeamRedPath    RedisPlayersPath = "teams.red.players"
	TeamBluePath   RedisPlayersPath = "teams.blue.players"
)

const (
	MsgJoinGame     MessageType = "join_game"
	MsgLeaveGame    MessageType = "leave_game"
	MsgGameState    MessageType = "game_state"
	MsgPlayerJoined MessageType = "player_joined"
	MsgPlayerLeft   MessageType = "player_left"
	MsgError        MessageType = "error"
)

type Message struct {
	Type   MessageType `json:"type"`
	Data   any         `json:"data"`
	GameID *uuid.UUID  `json:"game_id,omitempty"`
}

type Player struct {
	ID       uuid.UUID       `json:"id"`
	Name     string          `json:"name"`
	Conn     *websocket.Conn `json:"-"`
	GameID   uuid.UUID       `json:"-"`
	LastSeen time.Time       `json:"-"`
}

type Game struct {
	ID      uuid.UUID
	Players map[uuid.UUID]*Player
	mu      sync.RWMutex
	hub     *GameHub
}

func NewGame(id uuid.UUID, hub *GameHub) *Game {
	return &Game{
		ID:      id,
		Players: make(map[uuid.UUID]*Player),
		hub:     hub,
	}
}

func GameHubPlayerToGameStatePlayer(p *Player) dto.GameStatePlayer {
	return dto.GameStatePlayer{
		ID:   p.ID,
		Name: p.Name,
	}
}

func (g *Game) GetGameStateFromRedis(ctx context.Context) dto.GameState {
	gameKey := g.RedisGameKey()
	gameState, err := g.hub.rdb.JSONGet(ctx, gameKey, "$").Result()

	if err != nil {
		errStr := "Could not retrieve game state for game from redis"
		g.hub.logger.Error(errStr, "gameId", g.ID, "err", err)
		g.broadcastErrorMessage(ctx, errStr, err)
		panic(errStr)
	}

	var gs []dto.GameState
	err = json.Unmarshal([]byte(gameState), &gs)
	if err != nil || len(gs) < 1 {
		errStr := "Could not unmarshal game state"
		g.hub.logger.Error(errStr, "gameId", g.ID, "err", err)
		g.broadcastErrorMessage(ctx, errStr, err)
		panic(errStr)
	}
	return gs[0]
}

func (g *Game) broadcast(ctx context.Context, msg Message) {
	for _, player := range g.Players {
		if player.Conn != nil {
			if err := wsjson.Write(ctx, player.Conn, msg); err != nil {
				g.hub.logger.Error("Error broadcasting to player", "player", player.ID, "error", err)
				// Remove player if connection is broken
				go g.RemovePlayer(ctx, player.ID)
			}
		}
	}
}

func (g *Game) AddPlayer(ctx context.Context, player Player) {
	g.mu.Lock()
	defer g.mu.Unlock()
	player.GameID = g.ID
	player.LastSeen = time.Now()

	g.Players[player.ID] = &player

	err := g.AddPlayerToGameState(ctx, SpectatorsPath, &player)
	if err != nil {
		g.hub.logger.Error("Error adding player to reddis state", "err", err)
		g.broadcastErrorMessage(ctx, "Error adding player to reddis state", err)
		return
	}
	// Broadcast updated game state to all players in lobby
	g.broadcastGameState(ctx)
}

func (g *Game) broadcastGameState(ctx context.Context) {
	gameState := g.GetGameStateFromRedis(ctx)

	g.broadcast(ctx, Message{
		Type: MsgGameState,
		Data: gameState,
	})
}

func (g *Game) RemovePlayer(ctx context.Context, playerID uuid.UUID) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if player, exists := g.Players[playerID]; exists {
		// TODO: also remove from GS in redis
		// Handle disconnected state for players;
		// Disconnected field on player in redis state?
		delete(g.Players, playerID)

		// Notify other players that	 this player left
		// TODO: Send correct data
		g.broadcastGameState(ctx)
		// g.broadcast(ctx, Message{
		// 	Type: MsgPlayerLeft,
		// 	Data: "PlayerRemoved",
		// })

		// Close the connection
		if player.Conn != nil {
			player.Conn.Close(websocket.StatusNormalClosure, "Player Left the game")
		}

		// If lobby is empty, remove it
		if len(g.Players) == 0 {
			g.hub.RemoveGame(g.ID)
		}
	}
}

// AddPlayerToGameState checks in Redis with JSONPath, then adds
func (g *Game) AddPlayerToGameState(ctx context.Context, path RedisPlayersPath, player *Player) error {
	// TODO: add getter for viable paths if more than two teams
	paths := []RedisPlayersPath{SpectatorsPath, TeamBluePath, TeamRedPath}
	key := g.RedisGameKey()
	for _, p := range paths {
		// Query directly in Redis
		jsonPath := fmt.Sprintf("$.%s[?(@.id==\"%s\")]", p, player.ID)

		res, err := g.hub.rdb.JSONGet(ctx, key, jsonPath).Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if res != "[]" && res != "null" {
			// Do nothing cause player player already in lobby
			g.hub.logger.Info("Player already exists in", "player", player.ID, "path", p, "jsonGet", res)
			return nil
		}
	}

	// Append only if not found
	jsonPlayer, err := json.Marshal(GameHubPlayerToGameStatePlayer(player))
	if err != nil {
		return err
	}

	_, err = g.hub.rdb.JSONArrAppend(ctx, key, "$."+string(path), jsonPlayer).Result()

	return err
}

func (g *Game) broadcastErrorMessage(ctx context.Context, msg string, err error) {
	data := struct {
		Message string `json:"message"`
		Err     error  `json:"err"`
	}{Message: msg, Err: err}
	g.broadcast(ctx, Message{Type: MsgError, Data: data})
}

func (g *Game) RedisGameKey() string {
	return fmt.Sprintf("game:%s", g.ID)
}

type GameHub struct {
	games  map[uuid.UUID]*Game
	mu     sync.RWMutex
	logger *slog.Logger
	db     *database.DB
	rdb    *redis.Client
}

func NewGameHub(logger *slog.Logger, db *database.DB, rdb *redis.Client) *GameHub {
	return &GameHub{games: make(map[uuid.UUID]*Game), logger: logger, db: db, rdb: rdb}
}

func (h *GameHub) GetOrCreateGame(gameId uuid.UUID) *Game {
	h.mu.Lock()
	defer h.mu.Unlock()

	if game, exists := h.games[gameId]; exists {
		return game
	}

	game := NewGame(gameId, h)
	h.games[gameId] = game
	h.logger.Debug("Created new lobby:", "gameId", gameId)
	return game
}

func (gh *GameHub) GetGame(gameId uuid.UUID) *Game {
	gh.mu.RLock()
	defer gh.mu.RUnlock()
	return gh.games[gameId]
}

func (h *GameHub) RemoveGame(gameID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()
	// TODO: Delete game from redis as well
	if game, exists := h.games[gameID]; exists && len(game.Players) == 0 {
		delete(h.games, gameID)
		h.logger.Debug("Removed empty game", "gameId", gameID)
	}
}

func websocketPingLoop(ctx context.Context, c *websocket.Conn, userId uuid.UUID, gameId uuid.UUID, hub *GameHub) {
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
				hub.logger.Error("Failed ping for", "userId", userId.String())
				game := hub.GetGame(gameId)
				if game != nil {
					game.RemovePlayer(ctx, userId)
				}
				c.Close(websocket.StatusAbnormalClosure, "Failed ping response")
				return
			}
		}
	}
}

func processWSMessage(ctx context.Context, msg *Message, c *websocket.Conn, user sqlc.User, hub *GameHub) {
	switch msg.Type {
	case MsgJoinGame:
		game := hub.GetOrCreateGame(*msg.GameID)
		player := Player{ID: user.ID, Name: user.Name, Conn: c, GameID: *msg.GameID, LastSeen: time.Now()}
		game.AddPlayer(ctx, player)
	default:
		wsjson.Write(ctx, c, msg)
	}
}
