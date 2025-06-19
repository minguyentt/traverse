-- +goose Up
ALTER TABLE
    contracts
ADD
    COLUMN version INT DEFAULT 0;

-- +goose Down
ALTER TABLE
    contracts DROP COLUMN version;
