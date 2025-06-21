package storage

import (
	"context"
	"fmt"
	"traverse/internal/db"
	"traverse/models"

	"github.com/jackc/pgx/v5"
)

type ContractStorage interface {
	Create(ctx context.Context, contract *models.Contract) error
	Update(ctx context.Context, contract *models.Contract) error
	Delete(ctx context.Context, cID int64) error

	All(ctx context.Context, userID int64) ([]models.ContractMetaData, error)
	ByID(ctx context.Context, cID int64) (*models.Contract, error)
}

type contractStore struct {
	db *db.PGDB
}

func NewContractStore(db *db.PGDB) *contractStore {
	return &contractStore{db}
}

func (s *contractStore) Create(ctx context.Context, contract *models.Contract) error {
	txError := ExecTx(ctx, s.db, func(innerTx pgx.Tx) error {
		if err := s.create(ctx, contract, innerTx); err != nil {
			return err
		}

		// set the foreign key before inserting job details
		if contract.JobDetails != nil {
			contract.JobDetails.ContractID = contract.ID

			if err := s.insertJobDetailsWithContract(ctx, contract.JobDetails, innerTx); err != nil {
				return err
			}
		}

		return nil
	})

	if txError != nil {
		return fmt.Errorf("error occurred during transaction: %w", txError)
	}

	return nil
}

func (s *contractStore) Update(ctx context.Context, contract *models.Contract) error {
	txErr := ExecTx(ctx, s.db, func(innerTx pgx.Tx) error {
		if err := s.update(ctx, contract, innerTx); err != nil {
			return fmt.Errorf("failed to update contract: %w", err)
		}

		if err := s.updateJobDetailWithContract(ctx, contract.JobDetails, innerTx); err != nil {
			return fmt.Errorf("failed to upsert job details: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return fmt.Errorf("error occured during transaction: %w", txErr)
	}

	return nil
}

func (s *contractStore) Delete(ctx context.Context, cID int64) error {
	query := `
	DELETE FROM contracts
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	q, err := s.db.Exec(ctx, query, cID)
	if err != nil {
		return err
	}

	rows := q.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// retrieve contract by id
func (s *contractStore) ByID(ctx context.Context, cID int64) (*models.Contract, error) {
	query := `
		SELECT id, user_id, job_title, city, agency, created_at, updated_at, version
		FROM contracts
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var c models.Contract
	err := s.db.QueryRow(
		ctx,
		query,
		cID,
	).Scan(&c.ID, c.UserID, c.JobTitle, c.City, c.Agency, c.CreatedAt, c.UpdatedAt, c.Version)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	var cjd models.ContractJobDetails
	query = `
		SELECT contract_id, profession, assignment_length, experience
		FROM contract_job_details
		WHERE contract_id = $1
	`
	err = s.db.QueryRow(ctx, query, cID).
		Scan(&cjd.ContractID, &cjd.Profession, &cjd.AssignmentLength, &cjd.Experience)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to get job details: %w", err)
	}

	if err != pgx.ErrNoRows {
		c.JobDetails = &cjd
	}

	return &c, nil
}

// TODO: fetch the necessary data from contracts
// join the reviews and the review counts
func (s *contractStore) All(
	ctx context.Context,
	userID int64,
) ([]models.ContractMetaData, error) {
	// 1. grab the contract data
	// 2. get total reviews within the id of the contract
	// 3. left join the reviews table ON reviews.contract_id = contract.id
	// 4. left join the users ON contract.user_id = user.id
	query := `
	SELECT con.id, con.user_id, con.job_title, con.city, con.agency,
	u.id, u.firstname, u.username, u.email, con.version, con.created_at,
	COUNT(r.id) AS reviews_count
	FROM contracts con
	LEFT JOIN reviews r ON r.contract_id = con.id
	LEFT JOIN users u ON con.user_id = u.id
	GROUP BY con.id, u.id, u.username, con.created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying contracts: %w", err)
	}

	var contracts []models.ContractMetaData
	for rows.Next() {
		var c models.ContractMetaData
		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.JobTitle,
			&c.City,
			&c.Agency,
			&c.User.ID,
			&c.User.Firstname,
			&c.User.Username,
			&c.User.Email,
			&c.Version,
			&c.CreatedAt,
			&c.ReviewCounts,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning contract row: %w", err)
		}

		contracts = append(contracts, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating contract rows: %w", err)
	}

	return contracts, nil
}

func (s *contractStore) create(ctx context.Context, c *models.Contract, tx pgx.Tx) error {
	query := `
	INSERT INTO contracts (job_title, city, agency, user_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id, job_title, city, agency, user_id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRow(ctx, query, c.JobTitle, c.City, c.Agency, c.UserID).
		Scan(&c.ID, &c.JobTitle, &c.City, &c.Agency, &c.UserID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert contract: %w", err)
	}

	return nil
}

func (s *contractStore) insertJobDetailsWithContract(
	ctx context.Context,
	cjd *models.ContractJobDetails,
	tx pgx.Tx,
) error {
	query := `
	INSERT INTO contract_job_details (contract_id, profession, assignment_length, experience)
	VALUES ($1, $2, $3, $4)
	RETURNING contract_id, profession, assignment_length, experience
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRow(
		ctx,
		query,
		cjd.ContractID,
		cjd.Profession,
		cjd.AssignmentLength,
		cjd.Experience,
	).Scan(&cjd.ContractID, &cjd.Profession, &cjd.AssignmentLength, &cjd.Experience)
	if err != nil {
		return fmt.Errorf("failed to insert job details: %w", err)
	}

	return nil
}

func (s *contractStore) update(ctx context.Context, contract *models.Contract, tx pgx.Tx) error {
	query := `
	UPDATE contracts
	SET job_title = $1, city = $2, agency = $3, version = version + 1
	where id = $4 AND version = $5
	RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRow(ctx, query, contract.JobTitle, contract.City, contract.Agency, contract.ID).
		Scan(&contract.Version)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *contractStore) updateJobDetailWithContract(
	ctx context.Context,
	cjd *models.ContractJobDetails,
	tx pgx.Tx,
) error {
	/*
		1. detect any conflicts on the contract_id column due to unqiue constraint
		2. if a row is already inserted with an existing contract_id, trigger the conflict handler
		3. DO UPDATE SET to perform an UPDATE on the existing row instead of throwing an error
	*/
	query := `
        INSERT INTO contract_job_details (contract_id, profession, assignment_length, experience)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (contract_id)
        DO UPDATE SET
            profession = EXCLUDED.profession,
            assignment_length = EXCLUDED.assignment_length,
            experience = EXCLUDED.experience
		RETURNING contract_id, profession, assignment_length, experience
	`

	err := tx.QueryRow(
		ctx,
		query,
		cjd.ContractID,
		cjd.Profession,
		cjd.AssignmentLength,
		cjd.Experience,
	).Scan(&cjd.ContractID, &cjd.Profession, &cjd.AssignmentLength, &cjd.Experience)
	if err != nil {
		return fmt.Errorf("failed to upsert job details: %w", err)
	}

	return nil
}
