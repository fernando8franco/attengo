-- +goose Up
CREATE TABLE assistance_logs (
    id TEXT PRIMARY KEY,
    log_description TEXT NOT NULL,
    log_date TEXT NOT NULL DEFAULT (date('now')),
    entry_time TEXT DEFAULT (time('now')),
    exit_time TEXT,
    manual_minutes INTEGER NOT NULL DEFAULT 0,
    total_daily_minutes INTEGER GENERATED ALWAYS AS (
        IFNULL(
            ((strftime('%s', exit_time) - strftime('%s', entry_time)) / 60) + manual_minutes,
            0
        )
    ) VIRTUAL,
    user_id TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    deleted_at TEXT,
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE assistance_logs;