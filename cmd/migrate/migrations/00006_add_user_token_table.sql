-- +goose Up
CREATE TABLE IF NOT EXISTS user_tokens (
    user_id BIGINT NOT NULL,
    token bytea PRIMARY KEY
);

-- +goose Down
DROP TABLE IF EXISTS user_tokens;
