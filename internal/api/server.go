package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/minguyentt/traverse/configs"
	"github.com/minguyentt/traverse/internal/db"
	"go.uber.org/zap"
)

type server struct {
	ctx    context.Context
	db     *db.PGDB
	config *configs.Config
	mux    http.Handler

	logger *zap.SugaredLogger
}

func NewServer(
	ctx context.Context,
	db *db.PGDB,
	cfg *configs.Config,
	mux http.Handler,
	// service *services.Service,
	logger *zap.SugaredLogger,
) *server {
	return &server{
		ctx:    ctx,
		db:     db,
		config: cfg,
		mux:    mux,
		// service: service,
		logger: logger,
	}
}

func (s *server) Run() error {
	server := &http.Server{
		Addr:         ":" + s.config.SERVER.Port,
		Handler:      s.mux,
		ReadTimeout:  s.config.SERVER.ReadTimeout,
		WriteTimeout: s.config.SERVER.WriteTimeout,
	}

	if err := s.waitForPoolConn(); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	var wg sync.WaitGroup

	// FIX: tbh idk what these 2 routines are doing... lol
	go func() {
		s.logger.Infow("Server started listening on...", "addr", server.Addr)
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

	return nil
}

func (s *server) waitForPoolConn() error {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	// implement timer
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for database conection")
		case <-time.After(2 * time.Second):
			conn, err := s.db.GetConnection(ctx)
			if err == nil {
				conn.Release()
				s.logger.Info("succesfully connected to database")
				return nil
			}
			s.logger.Info("waiting for database connection...", "error", err)
		}
	}
}
