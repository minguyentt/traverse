package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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

	var userID int64

	// avoids precision loss in float64 type
	switch v := rawID.(type) {
	case float64:
		userID = int64(v)
	case string:
		userID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID format in string: %w", err)
		}
	default:
		return nil, fmt.Errorf("unexpected type for subject claim")
	}

	// if needed for downstream use
	claims["sub"] = strconv.FormatInt(userID, 10)

	return userID, nil
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
