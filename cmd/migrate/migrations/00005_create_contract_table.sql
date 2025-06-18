-- +goose Up
CREATE TABLE IF NOT EXISTS contracts
(
    id bigserial PRIMARY KEY,
    job_title text NOT NULL,
    city text NOT NULL,
    agency text NOT NULL,
    user_id bigint NOT NULL,
    created_at TIMESTAMP(0) NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS contract_job_details
(
    contract_id bigint NOT NULL,
    profession TEXT,
    assignment_length TEXT,
    experience TEXT,
    FOREIGN KEY (contract_id) REFERENCES contracts(id) ON DELETE CASCADE, -- when parent node deleted, all related childs are deleted as well
);

-- +goose Down
DROP TABLE IF EXISTS contracts;
DROP TABLE IF EXISTS contract_job_details;
