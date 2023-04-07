-- +goose Up
CREATE TABLE users_events (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    event_type TEXT NOT NULL,
    locked BOOL,
    deleted BOOL,
    payload jsonb,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP
);

-- +goose Down
DROP TABLE users_events;