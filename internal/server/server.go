package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
	"traverse/configs"
	"traverse/internal/db"
	"traverse/internal/routes"
)

type APIServer struct {
	ctx    context.Context
	pool   *db.Pool
	config *configs.Config
	router *routes.Router
	logger *slog.Logger
}

func NewApiServer(
	ctx context.Context,
	pool *db.Pool,
	cfg *configs.Config,
	r *routes.Router,
) *APIServer {
	logger := slog.Default().With("area", "API Server")

	return &APIServer{
		ctx:    ctx,
		pool:   pool,
		config: cfg,
		router: r,
		logger: logger,
	}
}

func (s *APIServer) Run() error {
    if err := s.waitForPoolConn(); err != nil {
        return fmt.Errorf("database connection failed: %w", err)
    }
	mux := s.router.SetupRouter()

	server := &http.Server{
		Addr:         ":" + s.config.SERVER.Port,
		Handler:      mux,
		ReadTimeout:  s.config.SERVER.ReadTimeout,
		WriteTimeout: s.config.SERVER.WriteTimeout,
	}

	var wg sync.WaitGroup

    //FIX: tbh idk what these 2 routines are doing... lol
	go func() {
		s.logger.Info("Server started", "Addr", server.Addr)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "server listening error: %s\n", err)
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		<-s.ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("Graceful shutdown error", "err", err)
		}
	}()
	wg.Wait()

	s.logger.Info("Server started and listening", "PORT", server.Addr)
    s.logger.Info("connected to the database")

	return nil
}

func (s *APIServer) waitForPoolConn() error {
    ctx, cancel := context.WithTimeout(s.ctx, 30 * time.Second)
    defer cancel()

    // implement timer
    for {
        select {
        case <-ctx.Done():
            return fmt.Errorf("timeout waiting for database conection")
        case <-time.After(2 * time.Second):
            conn, err := s.pool.GetConnection(ctx)
            if err == nil {
                conn.Release()
                s.logger.Info("succesfully connected to database")
                return nil
            }
            s.logger.Info("waiting for database connection...", "error", err)
        }
    }
}
