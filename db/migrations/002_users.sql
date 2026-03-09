-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL UNIQUE,
    required_hour_id INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TEXT,
    CONSTRAINT fk_required_hour
        FOREIGN KEY (required_hour_id)
        REFERENCES required_hours(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE users;