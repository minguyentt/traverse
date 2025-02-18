-- +goose Up
ALTER TABLE
    IF EXISTS users
ADD
    COLUMN account_type_id INT REFERENCES account_types(id) DEFAULT 1;

UPDATE
    users
SET account_type_id = (
    SELECT id
    FROM account_types
    WHERE _type = 'user'
);

ALTER TABLE users
ALTER COLUMN account_type_id DROP DEFAULT;

ALTER TABLE users
ALTER COLUMN account_type_id
SET NOT NULL;

-- +goose Down
