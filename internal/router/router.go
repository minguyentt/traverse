package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/minguyentt/traverse/internal/handlers"
	mw "github.com/minguyentt/traverse/internal/middlewares"
)

type Router struct {
	mux *chi.Mux
	middlewares *mw.Middlewares
}

func NewRouter(r *chi.Mux, middlewares *mw.Middlewares) *Router {
	return &Router{
        mux: r,
        middlewares: middlewares,
    }
}

func (r *Router) Mount(handlers *handlers.Handlers) http.Handler {
	r.mux.Use(r.middlewares.WithLogging())

	r.mux.Route("/v1", func(r chi.Router) {
		r.Get("/health", handlers.HealthChecker)

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Get("/", handlers.UserByID)
			})
		})
	})

	return r.mux
}
