-- +goose Up
ALTER TABLE
    contracts
ADD
    COLUMN updated_at timestamp(0) NOT NULL DEFAULT NOW();

-- +goose Down
ALTER TABLE
    contracts
DROP
    COLUMN updated_at;
