-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL UNIQUE,
    required_hour_id INTEGER NOT NULL,
    CONSTRAINT fk_required_hour
        FOREIGN KEY (required_hour_id)
        REFERENCES required_hours(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE users;