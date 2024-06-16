-- +goose Up
alter table users
add column api_key varchar(64) NOT NULL UNIQUE DEFAULT
encode(sha256(random()::text::bytea), 'hex');

-- +goose Down
ALTER TABLE users
DROP COLUMN api_key;
