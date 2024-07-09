-- +goose Up
CREATE TABLE feeds(id UUID PRIMARY KEY, name TEXT, url TEXT UNIQUE, user_id UUID REFERENCES users(id) ON DELETE CASCADE);

-- +goose Down
DROP TABLE feeds;