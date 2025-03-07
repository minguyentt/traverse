package storage

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"traverse/internal/db"
	"traverse/models"

	"github.com/jackc/pgx/v5"
)

type UserStorage interface {
	// user creation and retrieval
	CreateUser(ctx context.Context, user *models.User) error
	FetchAll(ctx context.Context) ([]models.User, error)
	Find(ctx context.Context, username string) (*models.User, error)
	ByID(ctx context.Context, userID int64) (*models.User, error)
	// FindByEmail(ctx context.Context, email string) (*models.User, error)

	// token management
	CreateTokenEntry(ctx context.Context, user_id int64, token string, exp time.Duration) error
	ActivateUserToken(ctx context.Context, token string) error

	DeleteUser(context.Context, int64) error
}

type userStore struct {
	db *db.PGDB
}

func NewUserStore(db *db.PGDB) *userStore {
	return &userStore{db}
}

func (s *userStore) FetchAll(ctx context.Context) ([]models.User, error) {
	query := `
    SELECT id, firstname, username, email, created_at
    FROM users
    `

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query row: %w", err)
	}
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var user models.User
        err := rows.Scan(&user.ID, &user.Firstname, &user.Username, &user.Email, &user.CreatedAt)
        if err != nil {
            return nil, fmt.Errorf("unable to scan row: %w", err)
        }
        users = append(users, user)
    }

    return users, nil
}

// executes db insertions to users & user_token tables
func (s *userStore) CreateUser(
	ctx context.Context,
	user *models.User,
) error {
	outer := ExecTx(ctx, s.db, func(inner pgx.Tx) error {
		if err := s.create(ctx, user, inner); err != nil {
			return err
		}

		return nil
	})

	if outer != nil {
		return outer
	}

	return nil
}

func (s *userStore) Find(ctx context.Context, username string) (*models.User, error) {
	q := `
    SELECT id, username, password, email, created_at
    FROM users
    WHERE username = $1
    `

	// var timeStamp time.Time
	user := &models.User{}
	err := s.db.QueryRow(ctx, q, username).
		Scan(&user.ID, &user.Username, &user.Password.Hash, &user.Email, &user.CreatedAt)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	// fmtStr := timeStamp.Format(time.RFC3339)
	// user.CreatedAt = fmtStr

	return user, nil
}

func (s *userStore) ByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `
    SELECT users.id, firstname, username, email, created_at
    FROM users
    WHERE users.id = $1
    `

	var user models.User
	// var timeStamp time.Time

	err := s.db.QueryRow(ctx, query, userID).
		Scan(&user.ID,
			&user.Firstname,
			&user.Username,
			&user.Email,
			&user.CreatedAt)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	// fmtStr := timeStamp.Format(time.RFC3339)
	// user.CreatedAt = fmtStr

	return &user, nil
}

// TODO: setup email retrieval and email trap logic
// func (s *userStore) FindByEmail(ctx context.Context, email string) (*models.User, error) {
//     return nil, nil
// }

func (s *userStore) CreateTokenEntry(
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
	slog.Info("user token entry creation executed", "output", cmd.String())

	return nil
}

// 1. Retrieve user by finding the token it belongs to by ID
// 2. clean up the user token after executed
func (s *userStore) ActivateUserToken(ctx context.Context, token string) error {
	return ExecTx(ctx, s.db, func(in pgx.Tx) error {
		user, err := s.findUserWithToken(ctx, token, in)
		if err != nil {
			return err
		}

		if err := s.update(ctx, user, in); err != nil {
			return err
		}

		if err := s.deleteUserToken(ctx, user.ID, in); err != nil {
			return err
		}

		return nil
	})
}

func (s *userStore) DeleteUser(ctx context.Context, userID int64) error {
	outer := ExecTx(ctx, s.db, func(inner pgx.Tx) error {
		if err := s.delete(ctx, userID, inner); err != nil {
			return err
		}
		return nil
	})

	if outer != nil {
		return outer
	}

	return nil
}

func (s *userStore) findUserWithToken(
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

	// var timestamp time.Time
	user := &models.User{}
	err := tx.QueryRow(ctx, que, token, time.Now()).
		Scan(&user.ID, &user.Firstname, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	// fmtStr := timestamp.Format(time.RFC3339)
	// user.CreatedAt = fmtStr

	return user, nil
}

func (s *userStore) create(ctx context.Context, user *models.User, tx pgx.Tx) error {
	query := `
    INSERT INTO users (firstname, username, password, email)
    VALUES ($1, $2, $3, $4)
    RETURNING id, created_at
    `

	// var timeStamp time.Time
	err := tx.QueryRow(ctx, query, user.Firstname, user.Username, user.Password.Hash, user.Email).
		Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	// fmt := timeStamp.Format(time.RFC3339)
	// user.CreatedAt = fmt

	return nil
}

func (s *userStore) update(ctx context.Context, user *models.User, tx pgx.Tx) error {
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

func (s *userStore) delete(ctx context.Context, userID int64, tx pgx.Tx) error {
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

// use this to clean up tokens to avoid dupes
func (s *userStore) deleteUserToken(ctx context.Context, user_id int64, tx pgx.Tx) error {
	q := `
    DELETE FROM user_tokens
    WHERE user_id = $1
    `

	cmd, err := tx.Exec(ctx, q, user_id)
	if err != nil {
		return err
	}

	slog.Info("deleted user token", "output", cmd.String())

	return nil
}
