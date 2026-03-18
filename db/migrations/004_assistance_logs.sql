-- +goose Up
CREATE TABLE assistance_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    log_description TEXT NOT NULL,
    log_date TEXT NOT NULL DEFAULT CURRENT_DATE,
    entry_time TEXT DEFAULT CURRENT_TIME,
    exit_time TEXT,
    manual_minutes INTEGER NOT NULL DEFAULT 0,
    total_daily_minutes INTEGER GENERATED ALWAYS AS (
        IFNULL(
            ((strftime('%s', exit_time) - strftime('%s', entry_time)) / 60) + manual_minutes,
            0
        )
    ) VIRTUAL,
    user_id INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TEXT,
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE assistance_logs;