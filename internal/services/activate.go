package services

import (
	"context"

	"traverse/internal/storage"
)

type ActivateService struct {
	store *storage.Storage
}

func NewActivationService(store *storage.Storage) *ActivateService {
	return &ActivateService{
		store: store,
	}
}

func (s *ActivateService) ActivateUser(ctx context.Context, token string) error {
    // set active
    err := s.store.Users.SetActive(ctx, token)
    if err != nil {
        return err
    }

    return nil
}
