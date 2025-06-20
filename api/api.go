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
	"traverse/internal/db/redis/cache"
	"traverse/internal/middlewares"
	"traverse/internal/services"
	"traverse/internal/storage"
	"traverse/pkg/validator"

	"github.com/redis/go-redis/v9"
)

// api version control
const version = "1.1.0"

type server struct {
	cfg    *configs.Config
	db     *db.PGDB
	redis  *redis.Client
	logger *slog.Logger
}

type api struct {
	*server
	ctx        context.Context
	mux        *router.Router
	handler    *handlers.Handlers
	middleware *middlewares.Middleware
	cache      cache.Cache
}

func New(
	cfg *configs.Config,
	db *db.PGDB,
	redis *redis.Client,
	logger *slog.Logger,
) *server {
	return &server{
		cfg:    cfg,
		db:     db,
		redis:  redis,
		logger: logger,
	}
}

func (s *server) SetupAPIV1(
	ctx context.Context,
	router *router.Router,
) (*api, error) {
	// wiring api deps

	rdsCache, err := cache.New(s.redis)
	assert.NoError(err, "failed to create new redis cache")

	jwt := auth.New(s.cfg.AUTH.Token.Secret, s.cfg.AUTH.Token.Aud, s.cfg.AUTH.Token.Iss)
	v := validator.NewValidator()
	storage := storage.New(s.db)
	service := services.New(storage, jwt)

	deps := handlers.HandlerDeps{
		Service:   service,
		Validator: v,
		Cache:     rdsCache,
	}

	handlers := handlers.New(&deps)

	mw := middlewares.New(jwt, service.Users, s.logger)

	assert.NotNil(handlers, "nil encounter")
	assert.NotNil(service, "nil encounter")
	assert.NotNil(storage, "nil encounter")

	// TODO: remove the unused dependencies
	api := &api{
		ctx:        ctx,
		server:     s,
		mux:        router,
		handler:    handlers,
		cache:      rdsCache,
		middleware: mw,
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
