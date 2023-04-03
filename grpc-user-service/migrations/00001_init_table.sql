-- +goose Up
CREATE TABLE users (
    id       BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    email    TEXT NOT NULL,
    name TEXT,
    surname TEXT,
    removed BOOL,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP
);

-- +goose Down
DROP TABLE users_events;