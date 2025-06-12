-- +goose Up
CREATE TABLE IF NOT EXISTS contracts
(
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    address text NOT NULL,
    city text NOT NULL,
    agency text NOT NULL,
    user_id bigint NOT NULL,
    created_at TIMESTAMP(0) NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS contracts;
