package storage

import (
	"context"
	"time"

	"github.com/minguyentt/traverse/internal/db"
	"github.com/minguyentt/traverse/internal/models"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	Users interface {
		CreateUser(ctx context.Context, user *models.User, token string, exp time.Duration) error
		UserByID(ctx context.Context, userID int64) (*models.User, error)
		SetActive(ctx context.Context, token string) error
		DeleteUser(context.Context, int64) error
	}
	AccountType interface {
		AccountAlias(context.Context, string) (*models.AccountType, error)
	}
}

func NewStorage(db *db.PGDB) *Storage {
	return &Storage{
		Users:       &UserStore{db},
		AccountType: &AccountTypeStore{db},
	}
}

func ExecTx(ctx context.Context, db *db.PGDB, fn func(pgx.Tx) error) error {
	outerTx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	// not sure if its proper way to pass ctx here...
	txWithCtx := context.Background()

	if err := fn(outerTx); err != nil {
		_ = outerTx.Rollback(txWithCtx)
		return err
	}

	return outerTx.Commit(txWithCtx)
}
