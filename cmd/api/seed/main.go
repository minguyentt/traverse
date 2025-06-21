package main

import (
	"context"
	"log/slog"
	"github.com/minguyentt/traverse/configs"
	"github.com/minguyentt/traverse/internal/assert"
	"github.com/minguyentt/traverse/internal/db"
	"github.com/minguyentt/traverse/internal/seed"
	"github.com/minguyentt/traverse/internal/storage"
)

func main() {
	ctx := context.Background()
	cfg := configs.Env
	logger := slog.Default()
	dbLogger := logger.With("area", "database pool connections for seeding")
	db, err := db.NewPoolConn(ctx, cfg.DB.String(), dbLogger)
	assert.NoError(err, "pool conn error", "msg", err)
	defer db.Close()

	store := storage.New(db)
	seed.Seed(store, db)
}
