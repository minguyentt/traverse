package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
	"traverse/api/errors"
	"traverse/configs"

	"github.com/golang-jwt/jwt/v5"
)

// response wrapper for log middleware
type response struct {
	http.ResponseWriter
	code int
}

// TokenAuthMiddleware validates JWT tokens and injects the authenticated user into the request context.
// It follows a fail-fast approach for security validations and provides detailed error responses.
// The middleware performs the following steps:
// 1. Validates the Authorization header format
// 2. Verifies the JWT token
// 3. Extracts and validates the user ID
// 4. Retrieves the user from storage
// 5. Injects the user into the request context
func (api *api) TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := parseAuthHeader(r, "Bearer")
		if err != nil {
			errors.UnauthorizedErr(w, r, err)
			return
		}

		claims, err := api.validateToken(token)
		if err != nil {
			errors.UnauthorizedErr(w, r, err)
			return
		}

		user, err := api.getUserFromClaims(r.Context(), claims)
		if err != nil {
			errors.UnauthorizedErr(w, r, err)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// example:
// basic Authentication string
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func (api *api) BasicAuthMiddleware(next http.Handler) http.Handler {
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

		api.logger.Info("Authorized admin user", "user", user)
		next.ServeHTTP(w, r)
	})
}

// middleware logger
func (api *api) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		resOut := &response{w, http.StatusOK}

		next.ServeHTTP(resOut, r)
		dur := time.Since(start)
		api.logger.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", resOut.code,
			"duration", dur,
			"ip", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
	})
}

// TODO: need rate limiter?

// parses the authorization request header with auth scheme type
// splits the header
// returns part[1]
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

// validate the jwt token and returns mapped claims
func (api *api) validateToken(token string) (jwt.MapClaims, error) {
	valid, err := api.jwt.Validate(token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	claims, ok := valid.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims format")
	}

	return claims, nil
}

// extract user id from claims and fetch user from storage
// TODO: fix this. i dont like it. why are we calling the storage here in the app layer?? the fuck.
func (api *api) getUserFromClaims(
	ctx context.Context,
	claims jwt.MapClaims,
) (any, error) {
	// extract subject claim
	rawID, ok := claims["sub"]
	if !ok {
		return nil, fmt.Errorf("missing subject claim")
	}

	userID, ok := rawID.(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID format")
	}

	user, err := api.storage.Users.ByID(ctx, int64(userID))
	if err != nil {
		return nil, fmt.Errorf("user lookup failed: %w", err)
	}

	return user, nil
}

// func (r *response) writeWriter(code int) {
// 	r.code = code
// 	r.ResponseWriter.WriteHeader(code)
// }
