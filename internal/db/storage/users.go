package storage

import (
	"context"
	"fmt"
	"traverse/internal/db"
	"traverse/internal/models"

	"github.com/jackc/pgx/v5"
)

type UserStore struct {
	db *db.PGDB
}

func (s *UserStore) CreateUser(ctx context.Context, user *models.User) error {
	// query := `
	// INSERT INTO users
	// `
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
		Scan(&user.Id, &user.Firstname, &user.Username, &user.Email, &user.IsActive, &user.AccountType, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("unable to query user id: %w", err)
	}

	return &user, nil
}
