package server

import (
	"expvar"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"traverse/configs"
	"traverse/api/handlers"
)

func (api *api) mount(handlers *handlers.Handlers) http.Handler {
	api.mux.Use(middleware.RequestID)
	api.mux.Use(middleware.RealIP)
	api.mux.Use(middleware.Logger)
	api.mux.Use(middleware.Recoverer)
	api.mux.Use(cors.Handler(configs.WithCorsOpts()))

	api.mux.Use(middleware.Timeout(60 * time.Second))
	// api.mux.Use(api.LoggerMiddleware)

	api.mux.Route("/v1", func(home chi.Router) {
        home.Route("/auth", func (auth chi.Router){
            // POST => registration handler
            // POST => token creation handler
        })

        // admin use routes
		home.Get("/health", handlers.HealthChecker)
		home.With(api.BasicAuthMiddleware).Get("/debug/vars", expvar.Handler().ServeHTTP)

		home.Route("/users", func(user chi.Router) {
            user.Put("/activate/{token}", handlers.ActivateUser)

			user.Route("/{userID}", func(r chi.Router) {
				r.Use(api.TokenAuthMiddleware)

				r.Get("/", handlers.UserByID)
			})
		})

	})

	return api.mux
}
