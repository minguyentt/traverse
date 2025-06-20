package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (auth *jwToken) Authenticate(ctx context.Context, r *http.Request) (any, error) {
	tokenStr, err := parseTokenFromHeader(r, "Bearer")
	if err != nil {
		return nil, err
	}

	valid, err := auth.Validate(tokenStr)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	claims, ok := valid.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims format")
	}

	// extract subject claim
	rawID, ok := claims["sub"]
	if !ok {
		return nil, fmt.Errorf("missing subject claim")
	}

	userID, ok := rawID.(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID format")
	}

	return int64(userID), nil
}

// responsible for fetching the header from user request
// splits the header parts and returns the bearer token
func parseTokenFromHeader(r *http.Request, authScheme string) (string, error) {
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
