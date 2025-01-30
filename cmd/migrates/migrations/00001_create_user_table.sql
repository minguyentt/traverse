-- +goose Up
CREATE TABLE IF NOT EXISTS users_test (
    id int NOT NULL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    user_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

INSERT INTO users_test VALUES
('1', 'maiko', 'nooyen', 'swag', NOW(), NOW()),
('2', 'luna', 'nooyen', 'kitty2', NOW(), NOW()),
('3', 'uwu', 'nooyen', 'kitty1', NOW(), NOW());
-- +goose Down
DROP TABLE IF EXISTS users_test;
