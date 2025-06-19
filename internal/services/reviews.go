package services

import (
	"context"
	"traverse/internal/storage"
	"traverse/models"
)

type ReviewService interface {
	GetReviewsByContractID(ctx context.Context, cID int64) ([]models.Review, error)
}

type reviewStore struct {
	store *storage.Storage
}

func NewReviewService(store *storage.Storage) *reviewStore {
	return &reviewStore{store}
}

func (s *reviewStore) GetReviewsByContractID(ctx context.Context, cID int64) ([]models.Review, error) {
	contractWithReviews, err := s.store.Reviews.GetByContractID(ctx, cID)
	if err != nil {
		return nil, err
	}

	return contractWithReviews, nil
}
