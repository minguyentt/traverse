package router

import (
	"github.com/go-chi/chi/v5"
)

type Router struct {
	*chi.Mux
}

func New() *Router {
    return &Router{chi.NewRouter()}
}
