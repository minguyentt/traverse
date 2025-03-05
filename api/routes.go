package server

import (
	"expvar"
	"net/http"
	"time"
	"traverse/configs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (api *api) mount() http.Handler {
	// api.mux.Use(middleware.RequestID)
	api.mux.Use(middleware.RealIP)
	api.mux.Use(middleware.Logger)
	api.mux.Use(middleware.Recoverer)
	api.mux.Use(cors.Handler(configs.WithCorsOpts()))

	api.mux.Use(middleware.Timeout(60 * time.Second))
	// api.mux.Use(api.LoggerMiddleware)


	api.mux.Route("/v1", func(r chi.Router) {
        r.Route("/register", func(pub chi.Router) {
            pub.Post("/user", api.handlers.Auth.RegisterUser)
        })

		// login/registration
		r.Route("/login", func(r chi.Router) {
			r.Post("/auth", api.handlers.Auth.Login)
		})

		// admin use routes
		r.Get("/health", api.handlers.HealthChecker)
		r.With(api.BasicAuthMiddleware).Get("/debug/vars", expvar.Handler().ServeHTTP)

		// users
		r.Route("/users", func(user chi.Router) {
			user.Route("/{userID}", func(r chi.Router) {
				r.Use(api.TokenAuthMiddleware)

				r.Get("/", api.handlers.Users.ByID)
			})
		})
	})

	return api.mux
}
