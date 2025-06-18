-- +goose Up
ALTER TABLE
    contracts
ADD
    CONSTRAINT fkey_contract_user FOREIGN KEY (user_id) REFERENCES users (id);

-- +goose Down
ALTER TABLE
    contracts
DROP
    CONSTRAINT fkey_contract_user;
