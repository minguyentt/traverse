package main

import (
	"context"
	"fmt"
	"time"
	"traverse/configs"
	"traverse/internal/server"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cfg := configs.LoadConfig()
	server := server.NewApiServer(ctx, cfg)
	err := server.Run()
	if err != nil {
		fmt.Printf("server error: %v", err)
	}
}
