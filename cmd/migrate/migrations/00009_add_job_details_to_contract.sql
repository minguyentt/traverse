-- +goose Up
CREATE TABLE IF NOT EXISTS contract_job_details
(
    id bigserial NOT NULL,
    contract_id bigint NOT NULL,
    profession TEXT NULL,
    assignment_length TEXT NULL,
    experience TEXT NULL
);

-- +goose Down
DROP TABLE IF EXISTS contract_job_details
