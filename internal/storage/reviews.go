package storage

import (
	"context"
	"traverse/internal/db"
	"traverse/models"
)

type ReviewStorage interface {
	Create(ctx context.Context, review *models.Review) error
	GetReviewsByContractID(ctx context.Context, cID int64) ([]models.Review, error)
}

type reviewStore struct {
	db *db.PGDB
}

func (s *reviewStore) Create(ctx context.Context, r *models.Review) error {
	query := `
	INSERT INTO reviews (contract_id, user_id, content)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRow(ctx, query, &r.ContractID, &r.UserID, &r.Content).Scan(&r.ID, &r.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *reviewStore) GetReviewsByContractID(
	ctx context.Context,
	cID int64,
) ([]models.Review, error) {
	return nil, nil
}
