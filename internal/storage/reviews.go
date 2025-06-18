package storage

import (
	"context"
	"traverse/internal/db"
	"traverse/models"
)

type ReviewStorage interface {
	Create(ctx context.Context, review *models.Review) error
	GetByContractID(ctx context.Context, cID int64) ([]models.Review, error)
}

type reviewStore struct {
	db *db.PGDB
}

func NewReviewStore(db *db.PGDB) *reviewStore {
	return &reviewStore{db}
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

func (s *reviewStore) GetByContractID(
	ctx context.Context,
	cID int64,
) ([]models.Review, error) {
	q := `
		SELECT rev.id, rev.contract_id, rev.user_id, u.username, u.id, rev.created_at
		FROM reviews rev
		JOIN users u ON u.id = rev.user_id
		WHERE rev.contract_id = $1
		ORDER BY rev.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.Query(ctx, q, cID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var r models.Review
		err := rows.Scan(&r.ID, &r.ContractID, &r.UserID, &r.Content, &r.User.Username, &r.User.ID)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}

	return reviews, nil
}
