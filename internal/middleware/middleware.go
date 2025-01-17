package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

const (
	Auth      = "auth"
	Logging   = "logging"
	RequestID = "request_id"
)

//TODO: make your own middleware ?
// come back later
type Middlewares map[string]func(http.Handler) http.Handler

func NewMiddlewares() Middlewares {
	return Middlewares{
		Logging:   middleware.Logger,
		RequestID: middleware.RequestID,
	}
}

