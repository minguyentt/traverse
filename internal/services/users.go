package services

import (
	"context"
	"traverse/api/models"
	"traverse/internal/auth"
	"traverse/internal/storage"
)

type UserService struct {
	store *storage.Storage
    authenticate auth.Authenticator
}

func NewUserService(store *storage.Storage, auth auth.Authenticator) *UserService {
	return &UserService{
		store: store,
        authenticate: auth,
	}
}

func (s *UserService) GetUser(ctx context.Context, username string) (*models.User, error) {
	user, err := s.store.Users.Retrieve(ctx, username)
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

func (s *UserService) RegisterUser(
	ctx context.Context,
	payload *models.RegistrationPayload,
) (*models.User, error) {
	// transfer obj
	user := &models.User{
		Firstname: payload.Firstname,
		Username:  payload.Username,
		Email:     payload.Email,
		AccountType: models.AccountType{
			AType: "user",
		},
	}

	// hash the password
	if err := user.Password.Hash([]byte(payload.Password)); err != nil {
		return nil, err
	}

	err := s.store.Users.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UserTokenEntry(
	ctx context.Context,
	payload *models.UserLoginPayload,
) (*models.UserToken, error) {
    user, err := s.store.Users.Retrieve(ctx, payload.Username)
    if err != nil {
        return nil, err
    }

    if err := user.Password.Compare([]byte(payload.Password)); err != nil {
        return nil, err
    }


    tokenStr, err := s.authenticate.CreateToken(user.ID, user.Username)
    if err != nil {
        return nil, err
    }

    userWithToken := &models.UserToken{
        User: user,
        Token: tokenStr,
    }

    //TODO: figure out when user is logged in. We need to create a token entry to the datbase...
    // either during login validation or during middlewares...

    return userWithToken, nil
}
