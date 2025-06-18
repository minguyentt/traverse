package main

import (
	"context"
	"log/slog"
	"traverse/configs"
	"traverse/internal/assert"
	"traverse/internal/db"
	"traverse/internal/seed"
	"traverse/internal/storage"
)

func main() {
	ctx := context.Background()
	cfg := configs.Env
	logger := slog.Default()
	dbLogger := logger.With("area", "database pool connections for seeding")
	db, err := db.NewPoolConn(ctx, cfg.DEV_DB.String(), dbLogger)
	assert.NoError(err, "pool conn error", "msg", err)
	defer db.Close()

	store := storage.New(db)
	seed.Seed(store, db)
}
