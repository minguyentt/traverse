-- +goose Up
CREATE TABLE IF NOT EXISTS account_types (
    id BIGSERIAL PRIMARY KEY,
    _type VARCHAR(255) NOT NULL UNIQUE,
    level INT NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO
  account_types (_type, description, level)
VALUES
  (
    'user',
    'A user can create contracts and post reviews',
    1
  );

INSERT INTO
  account_types (_type, description, level)
VALUES
  (
    'moderator',
    'A moderator can modify other users reviews',
    2
  );

INSERT INTO
  account_types (_type, description, level)
VALUES
  (
    'admin',
    'An admin can update/delete contracts and reviews',
    3
  );
-- +goose Down
DROP TABLE IF EXISTS account_types;
