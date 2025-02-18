-- +goose Up
ALTER TABLE
   user_tokens
ADD
    COLUMN expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL;

-- +goose Down
ALTER TABLE
    user_tokens DROP COLUMN expiry;
