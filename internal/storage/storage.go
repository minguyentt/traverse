package storage

import (
	"context"
	"errors"
	"time"

	"traverse/internal/db"
	"traverse/models"

	"github.com/jackc/pgx/v5"
)

var (
    ErrNotFound = errors.New("no resource found")
    ErrDuplicates = errors.New("found existing resource")
    ErrDuplicateUsername = errors.New("existing duplicate key for username")
)

type Storage struct {
	Users interface {
		CreateUser(ctx context.Context, user *models.User) error
        CreateTokenEntry(ctx context.Context, user_id int64, token string, exp time.Duration) error
        Find(ctx context.Context, username string) (*models.User, error)
		UserByID(ctx context.Context, userID int64) (*models.User, error)
		DeleteUser(context.Context, int64) error
	}
}

func NewStorage(db *db.PGDB) *Storage {
	return &Storage{
		Users:       &UserStore{db},
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
