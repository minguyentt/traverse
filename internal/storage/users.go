package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"time"
	"traverse/internal/db"
	"traverse/models"

	"github.com/jackc/pgx/v5"
)

type UserStore struct {
	db *db.PGDB
}

// executes db insertions to users & user_token tables
func (s *UserStore) CreateUser(
	ctx context.Context,
	user *models.User,
) error {
	outerTxErr := ExecTx(ctx, s.db, func(inner pgx.Tx) error {
		if err := s.create(ctx, user, inner); err != nil {
			return err
		}

		return nil
	})

	if outerTxErr != nil {
		return outerTxErr
	}

	return nil
}

func (s *UserStore) Find(ctx context.Context, username string) (*models.User, error) {
	q := `
    SELECT id, username, password, email, created_at
    FROM users
    WHERE username = $1
    `

	var timeStamp time.Time
	user := &models.User{}
	err := s.db.QueryRow(ctx, q, username).
		Scan(&user.ID, &user.Username, &user.Password.Hash, &user.Email, &timeStamp)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	fmtStr := timeStamp.Format(time.RFC3339)
	user.CreatedAt = fmtStr

	return user, nil
}

func (s *UserStore) CreateTokenEntry(
	ctx context.Context,
	user_id int64,
	token string,
	exp time.Duration,
) error {
	query := `
    INSERT INTO user_tokens (user_id, token, expiry)
    VALUES ($1, $2, $3)
    `

	cmd, err := s.db.Exec(ctx, query, user_id, token, time.Now().Add(exp))
	if err != nil {
		return err
	}
	// TOOD: remove later
	slog.Info("user token entry execution", "command", cmd.String())

	return nil
}

func (s *UserStore) UserByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `
    SELECT users.id, firstname, username, email, created_at
    FROM users
    WHERE users.id = $1
    `

	var user models.User
	var timeStamp time.Time

	err := s.db.QueryRow(ctx, query, userID).
		Scan(&user.ID,
			&user.Firstname,
			&user.Username,
			&user.Email,
			&timeStamp)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	fmtStr := timeStamp.Format(time.RFC3339)
	user.CreatedAt = fmtStr

	return &user, nil
}

func (s *UserStore) DeleteUser(ctx context.Context, userID int64) error {
	outerTxErr := ExecTx(ctx, s.db, func(inner pgx.Tx) error {
		if err := s.delete(ctx, userID, inner); err != nil {
			return err
		}
		return nil
	})

	if outerTxErr != nil {
		return outerTxErr
	}

	return nil
}

func (s *UserStore) create(ctx context.Context, user *models.User, tx pgx.Tx) error {
	query := `
    INSERT INTO users (firstname, username, password, email)
    VALUES ($1, $2, $3, $4)
    RETURNING id, created_at
    `

	var timeStamp time.Time
	err := tx.QueryRow(ctx, query, user.Firstname, user.Username, user.Password.Hash, user.Email).
		Scan(&user.ID, &timeStamp)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	fmt := timeStamp.Format(time.RFC3339)
	user.CreatedAt = fmt

	return nil
}

func (s *UserStore) update(ctx context.Context, user *models.User, tx pgx.Tx) error {
	query := `
    UPDATE users SET firstname = $1, username = $2, email = $3
    WHERE id = $5
    `

	_, err := tx.Exec(ctx, query, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) delete(ctx context.Context, userID int64, tx pgx.Tx) error {
	query := `
    SELECT id
    FROM users
    WHERE id = $1
    `

	_, err := tx.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) findUserWithToken(
	ctx context.Context,
	token string,
	tx pgx.Tx,
) (*models.User, error) {
	que := `
    SELECT id, firstname, username, email, created_at
    FROM users
    JOIN user_tokens ut ON id = ut.user_id
    WHERE ut.token = $1 AND ut.expiry > $2
    `

	chksum := sha256.Sum256([]byte(token))
	encodedToken := hex.EncodeToString(chksum[:])

	user := &models.User{}
	err := tx.QueryRow(ctx, que, encodedToken, time.Now()).
		Scan(&user.ID, &user.Firstname, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
