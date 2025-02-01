-- +goose Up
ALTER TABLE
    IF EXISTS users
ADD
    COLUMN account_type_id INT REFERENCES account_type(id) DEFAULT 1;

UPDATE
    users
SET account_type_id = (
    SELECT id FROM account_type
    WHERE alias = 'user'
);

-- TODO: WHAT IS THIS DOING? WTF
ALTER TABLE
    users
ALTER COLUMN
    account_type_id DROP DEFAULT;

ALTER TABLE
    users
ALTER COLUMN
    account_type_id
SET
    NOT NULL;

-- +goose Down
ALTER TABLE
    IF EXISTS users DROP COLUMN account_type_id;
