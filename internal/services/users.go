package services

import (
	"context"
	"log/slog"
	"time"
	cfg "traverse/configs"
	"traverse/internal/auth"
	"traverse/internal/storage"
	"traverse/models"

	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	store        *storage.Storage
	authenticate auth.Authenticator
}

func NewUserService(store *storage.Storage, auth auth.Authenticator) *UserService {
	return &UserService{
		store:        store,
		authenticate: auth,
	}
}

func (s *UserService) GetUser(ctx context.Context, username string) (*models.User, error) {
	user, err := s.store.Users.Find(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ByID(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.store.Users.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Register(
	ctx context.Context,
	payload *models.RegistrationPayload,
) (*models.User, error) {
	// transfer obj
	user := &models.User{
		Firstname: payload.Firstname,
		Username:  payload.Username,
		Email:     payload.Email,
	}

	// hash the password
	if err := user.Password.Set([]byte(payload.Password)); err != nil {
		return nil, err
	}

	err := s.store.Users.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(
	ctx context.Context,
	payload *models.UserLoginPayload,
) (*models.UserToken, error) {
	user, err := s.store.Users.Find(ctx, payload.Username)
	if err != nil {
		return nil, err
	}
	slog.Info("user login obj", "out", user)

	if err := user.Password.Compare([]byte(payload.Password)); err != nil {
		return nil, err
	}

	expiry := time.Now().Add(time.Hour * 24 * 3)

	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      expiry, // 3 day exp
		"iat":      time.Now().Unix(),
		"nbf":      time.Now().Unix(),
		"iss":      cfg.Env.AUTH.Token.Iss,
		"aud":      cfg.Env.AUTH.Token.Aud,
	}

	// generate the token and get the token str
	tokenStr, err := s.authenticate.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	// store the user token entry
	if err := s.store.Users.CreateTokenEntry(ctx, user.ID, tokenStr, expiry.Sub(time.Now())); err != nil {
		return nil, err
	}

	userWithToken := &models.UserToken{
		User:  user,
		Token: tokenStr,
	}

	return userWithToken, nil
}
