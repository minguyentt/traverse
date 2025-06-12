package storage

import (
	"context"
	"traverse/internal/db"
	"traverse/models"

	"github.com/jackc/pgx/v5"
)

type ContractStorage interface {
	Create(ctx context.Context, c *models.Contract) error
}

type contractStore struct {
	db *db.PGDB
}

func (s *contractStore) Create(ctx context.Context, c *models.Contract) error {
	query := `
	INSERT INTO contracts (title, address, city, agency, user_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRow(ctx, query, c.Title, c.Address, c.City, c.Agency, c.UserID).
		Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *contractStore) Update(ctx context.Context, c *models.Contract) error {
	query := `
	UPDATE contracts
	SET title = $1, address = $2, city = $3, agency = $4, version = version + 1
	where id = $5 AND version = $6
	RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRow(ctx, query, c.Title, c.Address, c.City, c.Agency, c.ID).Scan(&c.Version)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *contractStore) Delete(ctx context.Context, cID int64) error {
	query := `
	DELETE FROM contracts
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	q, err := s.db.Exec(ctx, query, cID)
	if err != nil {
		return err
	}

	rows := q.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// TODO: fetch the necessary data from contracts
// join the reviews and the review counts
func (s *contractStore) GetContractFeed(
	ctx context.Context,
	userID int64,
) ([]models.ContractMetaData, error) {
	// query := `
	// SELECT c.id, c.user_id, c.title, c.address, c.city, c.agency, u.username,
	// FROM contracts c
	// GROUP BY c.created_at
	// `

	var contracts []models.ContractMetaData

	return contracts, nil
}
