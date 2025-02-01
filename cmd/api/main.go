package main

import (
	"context"
	"fmt"
	"log"
	"traverse/configs"
	"traverse/internal/db"
	"traverse/internal/routes"
	"traverse/internal/server"
)

func main() {
	serverCtx := context.Background()
	poolCtx := context.Background()

	router := routes.NewRouter()

	cfg := configs.ENVS
	pool, err := db.PoolWithConfig(poolCtx, cfg.DB.String())
	if err != nil {
		log.Fatalf("server pool error: %v", err)
	}

	dsn := pool.Config()
	fmt.Println(dsn.ConnConfig.ConnString())

	defer pool.Close()

	server := server.NewApiServer(serverCtx, pool, cfg, router)
	err = server.Run()
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
