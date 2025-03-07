package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"
	cfg "traverse/configs"
	"traverse/internal/auth"
	"traverse/internal/storage"
	"traverse/models"

	"github.com/golang-jwt/jwt/v5"
)

type UserService interface {
    GetUsers(ctx context.Context) ([]models.User, error)
	GetUser(ctx context.Context, username string) (*models.User, error)
	UserByID(ctx context.Context, userID int64) (*models.User, error)

	RegisterUser(ctx context.Context, payload *models.RegistrationPayload) (*models.User, error)
	LoginUser(ctx context.Context, payload *models.UserLoginPayload) (*models.UserToken, error)

	ActivateUser(ctx context.Context, token string) error
}

type userService struct {
	store        *storage.Storage
	authenticate auth.Authenticator
}

func NewUserService(store *storage.Storage, auth auth.Authenticator) *userService {
	return &userService{
		store:        store,
		authenticate: auth,
	}
}

func (s *userService) GetUsers(ctx context.Context) ([]models.User, error) {
    users, err := s.store.Users.FetchAll(ctx)
    if err != nil {
        return nil, err
    }

    return users, nil
}

func (s *userService) GetUser(ctx context.Context, username string) (*models.User, error) {
	user, err := s.store.Users.Find(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UserByID(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.store.Users.ByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) RegisterUser(
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

func (s *userService) LoginUser(
	ctx context.Context,
	payload *models.UserLoginPayload,
) (*models.UserToken, error) {
	user, err := s.store.Users.Find(ctx, payload.Username)
	if err != nil {
		return nil, err
	}

	if err := user.Password.Compare([]byte(payload.Password)); err != nil {
		return nil, err
	}

	expiry := time.Hour * 24 * 3

	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(expiry).Unix(),
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

	chksum := sha256.Sum256([]byte(tokenStr))
	encodedToken := hex.EncodeToString(chksum[:])

	// store the user token entry
	if err := s.store.Users.CreateTokenEntry(ctx, user.ID, encodedToken, expiry); err != nil {
		return nil, err
	}

	userToken := &models.UserToken{
		Token:     tokenStr,
	}

	return userToken, nil
}

func (s *userService) ActivateUser(ctx context.Context, token string) error {
	if err := s.store.Users.ActivateUserToken(ctx, token); err != nil {
		return err
	}

	return nil
}
