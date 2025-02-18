package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/minguyentt/traverse/internal/models"
    "github.com/minguyentt/traverse/internal/db"
)

type AccountTypeStore struct {
	db *db.PGDB
}

// retrieving the account type from database:
// user (default), moderator, admin
func (s *AccountTypeStore) AccountAlias(
	ctx context.Context,
	alias string,
) (*models.AccountType, error) {
	query := `
    SELECT id, _type, level, description
    FROM account_type
    WHERE alias = $1
    `

	acc := &models.AccountType{}

	err := s.db.QueryRow(ctx, query, alias).
		Scan(&acc.ID, &acc.AType, &acc.Level, &acc.Description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("unable to query user id: %w", err)
	}

    return acc, nil
}
