package routes

import (
	"net/http"
	"traverse/internal/handlers"

	"github.com/go-chi/chi/v5"
)

// single route container
type route struct {
	method      string
	pattern     string
	handler     http.HandlerFunc
	description string
}

// route groups
type routeAPIVersion struct {
	pattern     string
	routes      []route
	description string
}

type Router struct {
	Routes []routeAPIVersion
	Handlers *handlers.APIHandlers
}

// TODO: should take in the config as argument later...
func NewRouter() *Router {
	return &Router{
		Routes: make([]routeAPIVersion, 0),
        Handlers: handlers.NewHandlers(),
	}
}

// its for using .Group() for chi router
func (r *Router) addAPIRoutes() []routeAPIVersion {
	return []routeAPIVersion{
		{
			pattern:     "/v1",
			description: "API v1 routes",
			routes: []route{
				{
					pattern: "/users",
					method:  http.MethodGet,
					handler: r.Handlers.Users.GetUsers,
					description: "List all Users",
				},
				{
					pattern:     "/contracts",
					method:      http.MethodGet,
					description: "List all Nurse Contracts",
				},
			},
		},
	}
}

func (r *Router) SetupRouter() *chi.Mux {
    mux := chi.NewRouter()
    groups := r.addAPIRoutes()

    for _, group := range groups {
        mux.Group(func (APIRouter chi.Router) {

            for _, route := range group.routes {
                routerWithRoutes := chi.NewRouter()

                routerWithRoutes.Method(route.method, "/", route.handler)
                APIRouter.Mount(group.pattern + route.pattern, routerWithRoutes)
            }
        })
    }

    return mux
}
