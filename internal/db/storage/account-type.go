package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"traverse/internal/db"
	"traverse/internal/models"
)

type AccountTypeStore struct {
	db *db.PGDB
}

func (s *AccountTypeStore) GetTypeByAlias(
	ctx context.Context,
	alias string,
) (*models.AccountType, error) {
	query := `
    SELECT id, alias, level, description
    FROM account_type
    WHERE alias = $1
    `

	acc := &models.AccountType{}

	err := s.db.QueryRow(ctx, query, alias).
		Scan(&acc.ID, &acc.Alias, &acc.Level, &acc.Description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("unable to query user id: %w", err)
	}

    return acc, nil
}
