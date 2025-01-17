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
	"traverse/internal/routes"

	"github.com/go-chi/chi/v5"
)

type APIServer struct {
	ctx    context.Context
	config *configs.Config
	logger *slog.Logger
}

func NewApiServer(ctx context.Context, cfg *configs.Config) *APIServer {
	logger := slog.Default().With("area", "API Server")

	return &APIServer{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}
}

// router tree nodes
func PrintTreeRouter(routes []chi.Route, indent string) {
	for _, route := range routes {
		fmt.Printf("%sPattern: %s\n", indent, route.Pattern)

		if len(route.Handlers) > 0 {
			fmt.Printf("%sHandlers:\n", indent)
			for method, handler := range route.Handlers {
				fmt.Printf("%s  %s -> %T\n", indent, method, handler)
			}
		}

		if route.SubRoutes != nil {
			fmt.Printf("%sSubroutes:\n", indent)
			PrintTreeRouter(route.SubRoutes.Routes(), indent+"  ")
		}
		fmt.Println()
	}
}

func (s *APIServer) Run() error {
    r := routes.NewRouter()
    mux := r.SetupRouter()

	PrintTreeRouter(mux.Routes(), "")

	server := &http.Server{
		Addr:         s.config.Port,
		Handler:      mux,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
	}

	var wg sync.WaitGroup

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

	return nil
}
