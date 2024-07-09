-- +goose Up
ALTER TABLE users ADD apikey VARCHAR(64) DEFAULT encode(sha256(random()::text::bytea), 'hex');

-- +goose Down
ALTER TABLE users DROP COLUMN apikey;