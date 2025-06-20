package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type TokenAuthenticator interface {
	GenerateToken(claims jwt.MapClaims) (string, error)
	Validate(token string) (*jwt.Token, error)
	Authenticate(ctx context.Context, r *http.Request) (any, error)
}

type jwToken struct {
	secretKey string
	audience  string
	issuer    string
}

func New(key, aud, iss string) *jwToken {
	return &jwToken{
		secretKey: key,
		audience:  aud,
		issuer:    iss,
	}
}

func (t *jwToken) GenerateToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (t *jwToken) Validate(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v\n", token.Header["alg"])
		}
		return []byte(t.secretKey), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(t.audience),
		jwt.WithIssuer(t.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
