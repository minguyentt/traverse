package api

import (
	"expvar"
	"net/http"
	"time"
	"github.com/minguyentt/traverse/configs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (api *api) mount() http.Handler {
	// api.mux.Use(middleware.RequestID)
	api.mux.Use(middleware.RealIP)
	// api.mux.Use(middleware.Logger)
	api.mux.Use(api.middleware.LoggerMiddleware)
	api.mux.Use(middleware.Recoverer)
	api.mux.Use(cors.Handler(configs.WithCorsOpts()))

	api.mux.Use(api.middleware.RateLimiterWithCMS)
	api.mux.Use(middleware.Timeout(60 * time.Second))

	api.mux.Route("/v1", func(r chi.Router) {
		api.mountUserRoutes(r)
		api.mountContractRoutes(r)
		api.mountAdminRoutes(r)
	})

	return api.mux
}

func (api *api) mountUserRoutes(r chi.Router) {
	r.Post("/register", api.handler.Registration)

	r.Route("/login", func(user chi.Router) {
		user.Post("/", api.handler.Login)
		user.Put("/activate", api.handler.ActivateUser)
	})

	r.Route("/users", func(user chi.Router) {
		user.Use(api.middleware.TokenAuth)
		user.Get("/{userID}", api.handler.GetUser)
	})
}

func (api *api) mountContractRoutes(r chi.Router) {
	r.Route("/contract", func(sub chi.Router) {
		sub.Use(api.middleware.TokenAuth)
		sub.Post("/", api.handler.CreateContract)

		sub.Route("/{id}", func(r chi.Router) {
			r.Use(api.middleware.TokenAuth)
			r.Use(api.middleware.LoadContractCtx)
			r.Get("/reviews", api.handler.ContractByID)
		})
	})

	r.Route("/feed", func(sub chi.Router) {
		sub.Use(api.middleware.TokenAuth)
		sub.Get("/", api.handler.Feed)
	})
}

func (api *api) mountAdminRoutes(r chi.Router) {
	r.Route("/system", func(sys chi.Router) {
		sys.Get("/health", api.handler.HealthChecker)

		sys.Group(func(admin chi.Router) {
			admin.Use(api.middleware.BasicAuth)
			// Or use JWT+role for better security in some cases:
			// admin.Use(api.TokenAuthMiddleware, api.AdminOnly)
			admin.Get("/debug/vars", expvar.Handler().ServeHTTP)
		})
	})
}
