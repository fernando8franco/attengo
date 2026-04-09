-- name: CreateEntryLog :one
INSERT INTO assistance_logs (id, log_description, user_id) 
VALUES (?, ?, ?)
RETURNING 
assistance_logs.id,
assistance_logs.entry_time,
assistance_logs.user_id,
(
    SELECT rh.total_minutes AS required_total
    FROM users u 
    JOIN required_hours rh ON u.required_hour_id = rh.id 
    WHERE u.id = assistance_logs.user_id
),
(
    SELECT CAST(SUM(total_daily_minutes) AS BIGINT) AS total_accumulated
    FROM assistance_logs al 
    WHERE al.user_id = assistance_logs.user_id
);

-- name: GetLastEntryLogByUser :one
SELECT id, log_date
FROM assistance_logs
WHERE user_id = ?
AND entry_time IS NOT NULL
AND exit_time IS NULL
ORDER BY log_date DESC
LIMIT 1;

-- name: UpdateExitLog :one
UPDATE assistance_logs 
SET exit_time = CURRENT_TIME
WHERE assistance_logs.id = ?
RETURNING 
assistance_logs.id,
assistance_logs.entry_time,
assistance_logs.exit_time,
assistance_logs.user_id,
(
    SELECT rh.total_minutes AS required_total
    FROM users u 
    JOIN required_hours rh ON u.required_hour_id = rh.id 
    WHERE u.id = assistance_logs.user_id
),
(
    SELECT CAST(SUM(total_daily_minutes) AS BIGINT) AS total_accumulated
    FROM assistance_logs al 
    WHERE al.user_id = assistance_logs.user_id
);

-- name: AddManualMinutes :one
INSERT INTO assistance_logs (id, log_description, user_id, manual_minutes, entry_time) 
values (?, ?, ?, ?, NULL)
RETURNING 
manual_minutes,
(
    SELECT name
    FROM users u  
    WHERE u.id = assistance_logs.user_id
),
(
    SELECT rh.total_minutes AS required_total
    FROM users u 
    JOIN required_hours rh ON u.required_hour_id = rh.id 
    WHERE u.id = assistance_logs.user_id
),
(
    SELECT CAST(SUM(total_daily_minutes) AS BIGINT) AS total_accumulated
    FROM assistance_logs al 
    WHERE al.user_id = assistance_logs.user_id
);