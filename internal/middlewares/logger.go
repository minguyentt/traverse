package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type response struct {
    http.ResponseWriter
    code int
}

type middlewareLogger struct {
	logger *zap.SugaredLogger
}

func NewMiddlewareLogger(logger *zap.SugaredLogger) *middlewareLogger {
    return &middlewareLogger{
        logger: logger,
    }
}

func (l *middlewareLogger) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

        outRes := l.response(w)

        next.ServeHTTP(outRes, r)

        dur := time.Since(start)

        l.logger.Infow("HTTP request",
            "method", r.Method,
            "path", r.URL.Path,
            "status", outRes.code,
            "duration", dur,
            "ip", r.RemoteAddr,
            "user_agent", r.UserAgent(),
            )
	})
}

func (l *middlewareLogger) response(w http.ResponseWriter) *response {
    return &response{
        ResponseWriter: w,
        code: http.StatusOK,
    }
}

func (r *response) write(code int) {
    r.code = code
    r.ResponseWriter.WriteHeader(code)
}
