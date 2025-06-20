package main

import (
	"context"
	"log/slog"
	"os"
	"time"
	"traverse/api/router"
	"traverse/configs"
	"traverse/internal/assert"
	"traverse/internal/db"
	"traverse/internal/db/redis"

	server "traverse/api"

	"github.com/lmittmann/tint"
	red "github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// loading configs
	cfg := configs.Env

	// TODO: move logger configurations somewhere else later..
	// currently writes all errors in red
	logger := slog.New(tint.NewHandler(os.Stderr, nil))
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))

	dbLogger := logger.With("area", "pgx")
	apiLogger := logger.With("area", "API")
	rdsLogger := logger.With("area", "redis")

	// setting up pool connections for db
	db, err := db.NewPoolConn(ctx, cfg.DEV_DB.String(), dbLogger)
	assert.NoError(err, "pool conn error", "msg", err)
	defer db.Close()

	// redis client
	var redisClient *red.Client
	if cfg.REDIS.Enabled {
		redisClient = redis.NewClient(cfg.REDIS.Addr, cfg.REDIS.Password, cfg.REDIS.DB)

		rdsLogger.Info("established connection with redis client")
		defer redisClient.Close()
	}

	router := router.New()

	server := server.New(cfg, db, redisClient, apiLogger)
	assert.NoError(err, "error setting up api routes", "err", err)

	v1API, err := server.SetupAPIV1(ctx, router)
	assert.NoError(err, "error setting up v1 api", "err", err)

	err = v1API.Run()
	assert.NoError(err, "api configuration error", "err", err)
	assert.NoError(err, "server error from starting", "msg", err)

	v1API.MonitorMetrics()
}
