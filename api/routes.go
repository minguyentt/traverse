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
	// api.mux.Use(middleware.Logger)
	api.mux.Use(api.LoggerMiddleware)
	api.mux.Use(middleware.Recoverer)
	api.mux.Use(cors.Handler(configs.WithCorsOpts()))

	api.mux.Use(middleware.Timeout(60 * time.Second))

	api.mux.Route("/v1", func(r chi.Router) {
		api.mountPublicRoutes(r)
		api.mountUserRoutes(r)
		api.mountContractRoutes(r)
		api.mountAdminRoutes(r)
	})

	return api.mux
}

func (api *api) mountPublicRoutes(r chi.Router) {
	r.Post("/register", api.handler.Registration)
	r.Post("/login", api.handler.Login)
}

func (api *api) mountUserRoutes(r chi.Router) {
	r.Route("/user", func(user chi.Router) {
		user.Group(func(g chi.Router) {
			g.Put("/activate/{token}", api.handler.ActivateUser)
		})

		user.Route("/{userID}", func(sub chi.Router) {
			sub.Use(api.TokenAuthMiddleware)
			sub.Get("/", api.handler.GetUser)
		})
	})

	// Admin Only routes
	r.Group(func(admin chi.Router) {
		admin.Route("/users", func(users chi.Router) {
			users.Get("/", api.handler.GetUsers)
		})
	})
}

func (api *api) mountContractRoutes(r chi.Router) {
	r.Route("/contracts", func(sub chi.Router) {
		sub.Use(api.TokenAuthMiddleware)
		sub.Post("/", api.handler.CreateContract)

		sub.Route("/contractID", func(r chi.Router) {
			r.Use(api.ContractMiddlewareCtx)
			r.Get("/", api.handler.ReviewsWithContractID)
		})
	})
}

func (api *api) mountAdminRoutes(r chi.Router) {
	r.Route("/system", func(sys chi.Router) {
		sys.Get("/health", api.handler.HealthChecker)

		sys.Group(func(admin chi.Router) {
			admin.Use(api.BasicAuthMiddleware)
			// Or use JWT+role for better security in some cases:
			// admin.Use(api.TokenAuthMiddleware, api.AdminOnly)
			admin.Get("/debug/vars", expvar.Handler().ServeHTTP)
		})
	})
}
