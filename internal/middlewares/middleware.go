package middlewares

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"traverse/configs"
	"traverse/internal/auth"
	"traverse/internal/ctx"
	"traverse/internal/services"
	"traverse/models"
	"traverse/pkg/errors"
)

type Middleware struct {
	// figure out what the mw deps are
	auth.TokenAuthenticator
	services.UserService
	logger *slog.Logger
}

func New(jwt auth.TokenAuthenticator, serv services.UserService, logger *slog.Logger) *Middleware {
	return &Middleware{
		jwt,
		serv,
		logger,
	}
}

func (m *Middleware) TokenAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := m.Authenticate(r.Context(), r)
		if err != nil {
			errors.UnauthorizedErr(w, r, err)
			return
		}

		// retrieve user data from db
		user, err := m.getUserFromRequest(r.Context(), userID.(int64))
		if err != nil {
			errors.InternalServerErr(w, r, err)
			return
		}

		// save the user context when user is logged in
		ctx := ctx.SetUser(r, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) getUserFromRequest(ctx context.Context, id int64) (*models.User, error) {
	user, err := m.UserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (m *Middleware) BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts, err := parseAuthHeader(r, "Basic")
		if err != nil {
			errors.UnauthorizedBasicErr(w, r, err)
			return
		}

		// decode the string as base64
		decode, err := base64.StdEncoding.DecodeString(parts)
		if err != nil {
			errors.UnauthorizedBasicErr(w, r, err)
		}

		// check creds from env variables
		user := configs.Env.AUTH.Admin.Username
		pass := configs.Env.AUTH.Admin.Password

		creds := strings.SplitN(string(decode), ":", 2)
		if len(creds) != 2 || creds[0] != user && creds[1] != pass {
			errors.UnauthorizedBasicErr(w, r, fmt.Errorf("invalid credentials"))
			return
		}

		m.logger.Info("Authorized admin user", "user", user)
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		resOut := &response{w, http.StatusOK}

		next.ServeHTTP(resOut, r)
		dur := time.Since(start)
		m.logger.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", resOut.code,
			"duration", dur,
			"ip", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
	})
}

// func (m *Middleware) RequireActivated(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		user, ok := r.Context().Value(auth.UserCtxKey).(*models.User)
// 		if !ok || user == nil || !user.Activated {
// 			http.Error(w, "account not activated", http.StatusForbidden)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

// response wrapper for log middleware
type response struct {
	http.ResponseWriter
	code int
}

func parseAuthHeader(r *http.Request, authScheme string) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", fmt.Errorf("authorization header is missing")
	}
	// parse and split header
	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != authScheme {
		return "", fmt.Errorf("malformed authorization header")
	}

	return parts[1], nil
}

// TODO: figure out how i want to utilize this redis cache for users...
// cache user middleware
