package main

import (
	"context"

	"github.com/minguyentt/traverse/configs"
	server "github.com/minguyentt/traverse/internal/api"
	"github.com/minguyentt/traverse/internal/db"
	"github.com/minguyentt/traverse/internal/handlers"
	"github.com/minguyentt/traverse/internal/middlewares"
	"github.com/minguyentt/traverse/internal/router"
	"github.com/minguyentt/traverse/internal/services"
	"github.com/minguyentt/traverse/internal/storage"
	"github.com/minguyentt/traverse/internal/zlogger"
)

func main() {
	ctx := context.Background()

	l := zlogger.NewLogger()
	dbLogger := l.WithArea("db connection")
	apiLogger := l.WithArea("api server")
	defer l.Sync()

	cfg := configs.ENVS
	db, err := db.NewPoolConn(ctx, cfg.DEV_DB.String(), dbLogger)
	if err != nil {
		l.Fatalf("server pool error: %v", err)
	}

	defer db.Close()

	storage := storage.NewStorage(db)
	service := services.NewServices(storage)
	handlers := handlers.NewHandlers(service)

	mw := middlewares.New(apiLogger)
	r := router.New(mw)

	mux := r.Mount(handlers)

	api := server.NewServer(ctx, db, cfg, mux, apiLogger)
	err = api.Run()
	if err != nil {
		l.Fatalf("server error: %v", err)
	}
}
