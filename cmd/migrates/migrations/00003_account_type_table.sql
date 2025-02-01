-- +goose Up
CREATE TABLE IF NOT EXISTS account_type (
    id BIGSERIAL PRIMARY KEY,
    alias VARCHAR(255) NOT NULL UNIQUE,
    level int NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO
    account_type (alias, level, description)
VALUES
('user', 1, 'User permission can create contracts, post reviews, moderate contract ownership'),
('moderator', 2, 'Moderator permission can modify user reviews'),
('admin', 3, 'Admin permission can modify ALL contracts and user reviews');

-- +goose Down
DROP TABLE IF EXISTS account_type;
