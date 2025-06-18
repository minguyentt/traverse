package storage

import (
	"context"
	"fmt"
	"traverse/internal/db"
	"traverse/models"

	"github.com/jackc/pgx/v5"
)

type ContractStorage interface {
	Create(ctx context.Context, c *models.Contract) error
	Update(ctx context.Context, c *models.Contract) error
	Delete(ctx context.Context, cID int64) error

	GetAllContracts(ctx context.Context, userID int64) ([]models.ContractMetaData, error)
	GetByID(ctx context.Context, cID int64) (*models.Contract, error)
}

type contractStore struct {
	db *db.PGDB
}

func (s *contractStore) Create(ctx context.Context, c *models.Contract) error {
	query := `
	INSERT INTO contracts (contract_name, address, city, agency, user_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRow(ctx, query, c.ContractName, c.Address, c.City, c.Agency, c.UserID).
		Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *contractStore) Update(ctx context.Context, c *models.Contract) error {
	query := `
	UPDATE contracts
	SET contract_name = $1, address = $2, city = $3, agency = $4, version = version + 1
	where id = $5 AND version = $6
	RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRow(ctx, query, c.ContractName, c.Address, c.City, c.Agency, c.ID).
		Scan(&c.Version)
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

// retrieve contract by id
func (s *contractStore) GetByID(ctx context.Context, cID int64) (*models.Contract, error) {
	query := `
	SELECT id, user_id, contract_name, address, city, agency, created_at, updated_at, version
	FROM contracts
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var c models.Contract
	err := s.db.QueryRow(
		ctx,
		query,
		cID,
	).Scan(&c.ID, c.UserID, c.ContractName, c.Address, c.City, c.Agency, c.CreatedAt, c.UpdatedAt, c.Version)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &c, nil
}

// TODO: fetch the necessary data from contracts
// join the reviews and the review counts
func (s *contractStore) GetAllContracts(
	ctx context.Context,
	userID int64,
) ([]models.ContractMetaData, error) {
	// 1. grab the contract data
	// 2. get total reviews within the id of the contract
	// 3. left join the reviews table ON reviews.contract_id = contract.id
	// 4. left join the users ON contract.user_id = user.id
	query := `
	SELECT con.id, con.user_id, con.contract_name, con.address, con.city, con.agency, u.username, con.version, con.created_at,
	COUNT(r.id) AS reviews_count
	FROM contracts c
	LEFT JOIN reviews r ON r.contract_id = con.id
	LEFT JOIN users u ON con.user_id = u.id
	GROUP BY p.id, u.username
	GROUP BY c.created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying contracts: %w", err)
	}

	var contracts []models.ContractMetaData
	for rows.Next() {
		var c models.ContractMetaData
		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.ContractName,
			&c.Address,
			&c.City,
			&c.Agency,
			&c.User.Username,
			&c.Version,
			&c.CreatedAt,
			&c.ReviewCounts,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning contract row: %w", err)
		}

		contracts = append(contracts, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating contract rows: %w", err)
	}

	return contracts, nil
}
