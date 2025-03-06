package storage

import (
	"context"
	"errors"
	"traverse/internal/db"

	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound          = errors.New("no resource found")
	ErrDuplicates        = errors.New("found existing resource")
	ErrDuplicateUsername = errors.New("existing duplicate key for username")
)

type Storage struct {
	Users UserStorage
}

func NewStorage(db *db.PGDB) *Storage {
	return &Storage{
		Users: &userStore{db},
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
