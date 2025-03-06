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
            pub.Post("/user", api.h.Auth.RegistrationHandler)
        })

		// login/registration
		r.Route("/login", func(r chi.Router) {
			r.Post("/", api.h.Auth.LoginHandler)
		})

		// admin use routes
		r.Get("/health", api.h.HealthChecker)
		r.With(api.BasicAuthMiddleware).Get("/debug/vars", expvar.Handler().ServeHTTP)

		// users
		r.Route("/users", func(user chi.Router) {
            user.Put("/activate/{token}", api.h.Auth.ActivationHandler)

			user.Route("/{userID}", func(r chi.Router) {
				r.Use(api.TokenAuthMiddleware)

				r.Get("/", api.h.Users.UserByIDHandler)
			})
		})
	})

	return api.mux
}
