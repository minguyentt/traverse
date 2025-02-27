package auth

import (
	"fmt"
	"time"
	cfg "traverse/configs"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator interface {
	CreateToken(userID int64, username string) (string, error)
	Validate(token string) (*jwt.Token, error)
}

type authenticator struct {
	secretKey string
	audience  string
	issuer    string
}

func NewJWTAuth(key, aud, iss string) *authenticator {
	return &authenticator{
		secretKey: key,
		audience:  aud,
		issuer:    iss,
	}
}

func (auth *authenticator) CreateToken(userID int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
        "username": username,
		"exp": time.Now().Add(time.Hour * 24 * 3), // 3 day exp
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": cfg.Env.AUTH.Token.Iss,
		"aud": cfg.Env.AUTH.Token.Aud,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(auth.secretKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (auth *authenticator) Validate(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v\n", token.Header["alg"])
		}
		return []byte(auth.secretKey), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(auth.audience),
		jwt.WithIssuer(auth.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
