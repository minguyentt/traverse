-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
    id BIGSERIAL PRIMARY KEY,
    firstname TEXT NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password BYTEA NOT NULL,
    email citext UNIQUE NOT NULL,
    created_at TIMESTAMP(0) NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;
