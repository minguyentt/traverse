-- +goose Up
CREATE TABLE IF NOT EXISTS reviews
(
   id BIGSERIAL PRIMARY KEY,
    contract_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP(0) NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS reviews;
