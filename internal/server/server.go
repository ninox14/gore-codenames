package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/lmittmann/tint"
	"github.com/ninox14/gore-codenames/internal/database"
	"github.com/ninox14/gore-codenames/internal/env"
)

type config struct {
	baseURL  string
	httpPort int
	cookie   struct {
		secretKey string
	}
	jwt struct {
		secretKey string
	}
}

type Server struct {
	port   int
	logger *slog.Logger
	db     database.Service
	config config
}

func NewServer() *http.Server {

	var cfg config

	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:8080")
	cfg.httpPort = env.GetInt("PORT", 8080)
	cfg.cookie.secretKey = env.GetString("COOKIE_SECRET_KEY", "d4q4sl5zd3exvpfnn5eu776ghd4up2z6")
	cfg.jwt.secretKey = env.GetString("JWT_SECRET_KEY", "5il7lpknmngmaklaquxzzfz7x5on3pxf")

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))
	NewServer := &Server{
		port:   cfg.httpPort,
		logger: logger,
		config: cfg,

		db: database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
