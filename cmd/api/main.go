package main

import (
	"context"
	"log/slog"
	"traverse/api/handlers"
	"traverse/api/router"
	"traverse/configs"
	"traverse/internal/assert"
	"traverse/internal/auth"
	"traverse/internal/db"
	"traverse/internal/services"
	"traverse/internal/storage"

	server "traverse/api"
)

func main() {
	ctx := context.Background()

	// loading configs
	cfg := configs.Env

	// setup loggers
	logger := slog.Default()
	dbLogger := logger.With("area", "database pool connections")
	sl := logger.With("area", "API Server")

	// setting up pool connections for db
	db, err := db.NewPoolConn(ctx, cfg.DEV_DB.String(), dbLogger)
	assert.NoError(err, "pool conn error", "msg", err)
	defer db.Close()

    // jwt setup
	jwt := auth.NewJWTAuth(cfg.AUTH.Token.Secret, cfg.AUTH.Token.Aud, cfg.AUTH.Token.Iss)

	// routes
	router := router.New()
	storage := storage.NewStorage(db)
	service := services.NewServices(storage, jwt)
	handlers := handlers.NewHandlers(service)

	server := server.NewServer(db, cfg, sl)
	api, err := server.SetupAPI(ctx, server, router, handlers, service, storage)
	assert.NoError(err, "error setting up api routes", "err", err)

	err = api.Run()
	assert.NoError(err, "api configuration error", "err", err)
	assert.NoError(err, "server error from starting", "msg", err)

	server.MonitorMetrics()
}
