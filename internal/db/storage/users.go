package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"traverse/internal/db"
	"traverse/internal/models"
)

type UserStore struct {
	db *db.PGDB
}

func (s *UserStore) CreateUser(ctx context.Context, user *models.User) error {
	query := `
	INSERT INTO users (id, username, password, email, account_type_id)
    VALUES ($1, $2, $3, (SELECT id FROM account_type WHERE alias = $4))
    RETURNING id, username, created_at
	`

	accAlias := user.AccountType.Alias
	if accAlias == "" {
		accAlias = "user"
	}

	// NOTE: probably when we begin the transaction
	// passing the ctx should contain validation

	tx, err := s.db.BeginWithTx(ctx)
	if err != nil {
		return err
	}

	err = tx.QueryRow(ctx, query, user.Username, user.Password, user.Email, accAlias).
		Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `
    SELECT users.id, firstname, username, email, created_at, account_status.*
    FROM users
    JOIN account_status ON (users.account_type = account_status.account_type)
    WHERE users.id = $1 AND is_active = true
    `

	var user models.User

	err := s.db.QueryRow(ctx, query, userID).
		Scan(&user.ID, &user.Firstname, &user.Username, &user.Email, &user.IsActive, &user.AccountType, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("unable to query user id: %w", err)
	}

	return &user, nil
}

func (s *UserStore) DeleteUser(ctx context.Context, userID int64) error {
	query := `
    SELECT id
    FROM users
    WHERE id = $1
    `

	tx, err := s.db.BeginWithTx(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, userID)
    if err != nil {
        return err
    }

    return nil

}
