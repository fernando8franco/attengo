-- +goose Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    is_admin BOOLEAN NOT NULL DEFAULT 0,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    required_hour_id INTEGER,
    period_id INTEGER,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TEXT,
    CONSTRAINT fk_required_hour
        FOREIGN KEY (required_hour_id)
        REFERENCES required_hours(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_period_id
        FOREIGN KEY (period_id)
        REFERENCES periods(id)
        ON DELETE CASCADE,
    CONSTRAINT uc_email_required_hour_id UNIQUE (email, required_hour_id, period_id)
);

-- +goose Down
DROP TABLE users;