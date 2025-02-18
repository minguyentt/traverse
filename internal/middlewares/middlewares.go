package middlewares

import (
	"net/http"

	"go.uber.org/zap"
)

type Middlewares struct {
	*middlewareLogger
}

func New(zlog *zap.SugaredLogger) *Middlewares {
	return &Middlewares{
		NewMiddlewareLogger(zlog),
	}
}

func (m *Middlewares) WithLogging() func(http.Handler) http.Handler {
	return m.Logger
}
