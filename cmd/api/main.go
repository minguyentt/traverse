package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/minguyentt/traverse/api"
	"github.com/minguyentt/traverse/api/router"
	"github.com/minguyentt/traverse/configs"
	"github.com/minguyentt/traverse/internal/assert"
	"github.com/minguyentt/traverse/internal/db"
	"github.com/minguyentt/traverse/internal/db/redis"

	"github.com/lmittmann/tint"
	red "github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	assert.NoError(err, "failed to load environment variables", "msg", err)
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
	apiLogger := logger.WithGroup("API")
	rdsLogger := logger.With("area", "redis")

	// setting up pool connections for db
	db, err := db.NewPoolConn(ctx, cfg.DB.String(), dbLogger)
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

	server := api.New(cfg, db, redisClient, apiLogger)
	assert.NoError(err, "error setting up api routes", "err", err)

	v1API, err := server.SetupAPIV1(ctx, router)
	assert.NoError(err, "error setting up v1 api", "err", err)

	err = v1API.Run()
	assert.NoError(err, "api configuration error", "err", err)
	assert.NoError(err, "server error from starting", "msg", err)

	v1API.MonitorMetrics()
}
