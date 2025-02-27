package server

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"traverse/configs"
	"traverse/internal/auth"
	"traverse/internal/db"
	"traverse/api/handlers"
	"traverse/api/router"
	"traverse/internal/services"
	"traverse/internal/storage"
)

// api version control
const version = "1.1.0"

type server struct {
	cfg    *configs.Config
	db     *db.PGDB
	logger *slog.Logger
}

type api struct {
	ctx      context.Context
	mux      *router.Router
	handlers *handlers.Handlers
	service  *services.Service
	storage  *storage.Storage
	*server
}

func NewServer(
	db *db.PGDB,
	cfg *configs.Config,
	logger *slog.Logger,
) *server {
	return &server{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *server) SetupAPI(
	ctx context.Context,
	server *server,
	router *router.Router,
	handlers *handlers.Handlers,
	service *services.Service,
	storage *storage.Storage,
) (*api, error) {
	api := &api{
		server:   server,
		mux:      router,
		handlers: handlers,
		service:  service,
		storage:  storage,
	}

	return api, nil
}

func (api *api) Run() error {
	mux := api.mount(api.handlers)

	if err := api.waitConnection(); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	server := &http.Server{
		Addr:         ":" + api.cfg.SERVER.Port,
		Handler:      mux,
		ReadTimeout:  api.cfg.SERVER.ReadTimeout,
		WriteTimeout: api.server.cfg.SERVER.WriteTimeout,
	}

	serverSig := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		api.logger.Info("server signal found", "msg", sig)

		serverSig <- server.Shutdown(ctx)
	}()

	api.logger.Info("Server started listening on...", slog.String("addr", server.Addr))
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-serverSig
	if err != nil {
		return err
	}

	return nil
}

func (api *api) waitConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// implement timer
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for database conection")
		case <-time.After(2 * time.Second):
			conn, err := api.server.db.GetConnection(ctx)
			if err == nil {
				conn.Release()
				api.logger.Info("succesfully connected to database")
				return nil
			}
			api.logger.Info("waiting for database connection...", "error", err)
		}
	}
}

func (s *server) MonitorMetrics() {
	expvar.NewString("version").Set(version)
	expvar.Publish("database connection pooling", expvar.Func(func() any {
		return s.db.Stat()
	}))

	expvar.Publish("concurrency", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
}
