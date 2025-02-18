package services

import (
	"context"

	"github.com/minguyentt/traverse/internal/models"
	"github.com/minguyentt/traverse/internal/storage"
)

type UserService struct {
    store *storage.Storage
}

func NewUserService(store *storage.Storage) *UserService {
	return &UserService{
        store: store,
	}
}

func (s *UserService) ByID(ctx context.Context, userID int64) (*models.User, error) {
    user, err := s.store.Users.UserByID(ctx, userID)
    if err != nil {
        return nil, err
    }

    return user, nil
}

