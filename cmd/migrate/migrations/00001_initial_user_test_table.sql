-- +goose Up
CREATE TABLE IF NOT EXISTS users_test (
    id int NOT NULL PRIMARY KEY,
    firstname TEXT NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

INSERT INTO users_test VALUES
('1', 'maiko', 'freshoe', '$2a$10$EixZaXNyZWR0aXJlMTIzNDU2c3RpbmRvbWV0aG9uZ2F0ZQ==', NOW(), NOW()),
('2', 'luna', 'luna_test', 'cutecat1', NOW(), NOW()),
('3', 'uwu', 'uwu_test', 'cutecat2', NOW(), NOW());
-- +goose Down
DROP TABLE IF EXISTS users_test;
