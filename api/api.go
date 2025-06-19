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
	"traverse/api/handlers"
	"traverse/api/router"
	"traverse/configs"
	"traverse/internal/assert"
	"traverse/internal/auth"
	"traverse/internal/db"
	"traverse/internal/services"
	"traverse/internal/storage"
	json "traverse/pkg/validator"

	"github.com/go-playground/validator/v10"
)

// api version control
const version = "1.1.0"

type server struct {
	cfg    *configs.Config
	db     *db.PGDB
	logger *slog.Logger
}

type api struct {
	*server
	ctx     context.Context
	mux     *router.Router
	handler *handlers.Handlers
	service *services.Service
	storage *storage.Storage

	jwt       auth.Authenticator
	validator *validator.Validate
}

func New(
	cfg *configs.Config,
	db *db.PGDB,
	logger *slog.Logger,
) *server {
	return &server{
		cfg:    cfg,
		db:     db,
		logger: logger,
	}
}

func (s *server) SetupAPIV1(
	ctx context.Context,
	router *router.Router,
) (*api, error) {
	jwtAuth := auth.NewJWTAuth(s.cfg.AUTH.Token.Secret, s.cfg.AUTH.Token.Aud, s.cfg.AUTH.Token.Iss)
	validator := json.NewValidator()
	storage := storage.New(s.db)
	service := services.New(storage, jwtAuth)
	handlers := handlers.New(service, validator)
	assert.NotNil(handlers, "nil encounter")
	assert.NotNil(service, "nil encounter")
	assert.NotNil(storage, "nil encounter")

	api := &api{
		ctx:       ctx,
		server:    s,
		mux:       router,
		handler:   handlers,
		service:   service,
		storage:   storage,
		jwt:       jwtAuth,
		validator: validator,
	}

	return api, nil
}

func (api *api) Run() error {
	mux := api.mount()

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
	api.logger.Info("API Ready. Waiting for requests...")

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
			api.logger.Warn("waiting for database connection...", "msg", err)
		}
	}
}

func (s *api) MonitorMetrics() {
	expvar.NewString("version").Set(version)
	expvar.Publish("database connection pooling", expvar.Func(func() any {
		return s.db.Stat()
	}))

	expvar.Publish("concurrency", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
}
