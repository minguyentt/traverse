package services

import (
	"context"
	"github.com/minguyentt/traverse/models"
)

func (s *contractService) ReviewsWithContractID(ctx context.Context, cID int64) ([]models.Review, error) {
	contractWithReviews, err := s.store.Reviews.GetByContractID(ctx, cID)
	if err != nil {
		return nil, err
	}

	return contractWithReviews, nil
}
