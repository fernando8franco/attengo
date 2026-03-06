CREATE TABLE assistance_logs (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    required_hour_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_required_hour
        FOREIGN KEY (required_hour_id)
        REFERENCES required_hours(id)
        ON DELETE CASCADE
);