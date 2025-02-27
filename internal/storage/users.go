package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"time"
	"traverse/internal/db"
	"traverse/api/models"

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

// retrieve user by username
func (s *UserStore) Retrieve(ctx context.Context, username string) (*models.User, error) {
	q := `
    SELECT id, username, password, email, created_at
    FROM users
    WHERE username = $1 AND is_active = true
    RETURNING id
    `

	user := &models.User{}
	err := s.db.QueryRow(ctx, q, username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt)
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

func (s *UserStore) UserTokenEntry(
    ctx context.Context,
    user_id int64,
    token string,
    exp time.Duration,
) error {
    query := `
    INSERT INTO user_tokens (user_id, token, exp)
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

func (s *UserStore) update(ctx context.Context, user *models.User, tx pgx.Tx) error {
    query := `
    UPDATE users SET firstname = $1, username = $2, email = $3, is_active $4
    WHERE id = $5
    `

    _, err := tx.Exec(ctx, query, user.ID)
    if err != nil {
        return err
    }

    return nil
}

func (s *UserStore) UserByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `
    SELECT users.id, firstname, username, email, is_active, account_types.*, created_at
    FROM users
    JOIN account_types ON (users.account_type_id = account_types.id)
    WHERE users.id = $1 AND is_active = true
    `

	var user models.User
	var timeStamp time.Time

	err := s.db.QueryRow(ctx, query, userID).
		Scan(&user.ID,
			&user.Firstname,
			&user.Username,
			&user.Email,
			&user.IsActive,
			&user.AccountType.ID,
			&user.AccountType.AType,
			&user.AccountType.Level,
			&user.AccountType.Description,
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

func (s *UserStore) SetActive(ctx context.Context, token string) error {
	outerErr := ExecTx(ctx, s.db, func(inner pgx.Tx) error {
		user, err := s.findUserWithToken(ctx, token, inner)
		if err != nil {
			return err
		}

		user.IsActive = true
		if err := s.update(ctx, user, inner); err != nil {
			return err
		}

		return nil
	})

	if outerErr != nil {
		return outerErr
	}

	return nil
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
    INSERT INTO users (id, username, password, email, account_type_id)
    VALUES ($1, $2, $3, (SELECT id FROM account_types WHERE _type = $4))
    RETURNING id, username, created_at
    `

	accType := user.AccountType.AType
	if accType == "" {
		accType = "user"
	}

	err := tx.QueryRow(ctx, query, user.Username, user.Password, user.Email, accType).
		Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}


func (s *UserStore) findUserWithToken(
	ctx context.Context,
	token string,
	tx pgx.Tx,
) (*models.User, error) {
	que := `
    SELECT id, firstname, username, email, created_at, is_active
    FROM users
    JOIN user_tokens ut ON id = ut.user_id
    WHERE ut.token = $1 AND ut.expiry > $2
    `

	chksum := sha256.Sum256([]byte(token))
	encodedToken := hex.EncodeToString(chksum[:])

	user := &models.User{}
	err := tx.QueryRow(ctx, que, encodedToken, time.Now()).
		Scan(&user.ID, &user.Firstname, &user.Username, &user.Email, &user.CreatedAt, &user.CreatedAt)
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
