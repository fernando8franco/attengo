-- name: CreateEntryLog :one
INSERT INTO assistance_logs (id, log_description, user_id) 
VALUES (?, ?, ?)
RETURNING 
id,
entry_time,
exit_time,
total_daily_minutes,
user_id;

-- name: GetIDFromLastEntryLogByUser :one
SELECT id
FROM assistance_logs
WHERE user_id = ?
AND exit_time IS NULL
ORDER BY log_date DESC
LIMIT 1;

-- name: UpdateEntryLog :one
UPDATE assistance_logs 
SET exit_time = CURRENT_TIME
WHERE id = ?
RETURNING
id,
entry_time,
exit_time,
total_daily_minutes,
user_id;