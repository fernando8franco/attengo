-- name: CreateUser :one
INSERT INTO users (id, name, email, password, required_hour_id, period_id) 
VALUES (?, ?, ?, ?, ?, ?)
RETURNING 
users.id,
users.name,
users.email,
users.password,
(
  SELECT required_hours.type
  FROM required_hours
  WHERE users.required_hour_id = required_hours.id
),
(
  SELECT required_hours.total_minutes AS hours
  FROM required_hours
  WHERE users.required_hour_id = required_hours.id
),
(
  SELECT periods.name AS period
  FROM periods
  WHERE users.period_id = periods.id
),
(
  SELECT CAST(COALESCE(SUM(total_daily_minutes), 0) AS BIGINT) AS total_hours
  FROM assistance_logs
  WHERE users.id = assistance_logs.user_id
);

-- name: GetUsersPasswords :many
SELECT password FROM users WHERE is_admin = 0;

-- name: ValidateUserPassword :one
SELECT id FROM users
WHERE is_admin = 0 AND password = ?;

-- name: CreateAdmin :one
INSERT INTO users (id, is_admin, name, email, password) 
VALUES (?, 1, ?, ?, ?)
RETURNING 
id;

-- name: ExistsAdmin :one
SELECT EXISTS (
  SELECT 1
  FROM users
  WHERE is_admin = 1
) = 1;

-- name: GetAdminIDAndPasswordByEmail :one
SELECT id, password
FROM users
WHERE email = ?
AND is_admin = 1
AND deleted_at IS NULL;

-- name: GetActiveUsers :many
SELECT DISTINCT users.id, users.name, CAST(time(assistance_logs.entry_time, '-6 hours') AS TEXT) AS entry_time
FROM users
INNER JOIN assistance_logs ON users.id = assistance_logs.user_id
WHERE users.is_admin = 0 
AND users.deleted_at IS NULL
AND assistance_logs.entry_time IS NOT NULL
AND assistance_logs.exit_time IS NULL
AND assistance_logs.log_date = date('now', '-6 hours');

-- name: GetNotAdminUsers :many
SELECT 
    users.id, 
    users.name, 
    users.email, 
    periods.name AS period, 
    required_hours.type, 
    required_hours.total_minutes AS hours,
    CAST(COALESCE(SUM(assistance_logs.total_daily_minutes), 0) AS BIGINT) AS total_hours,
    users.password
FROM users
INNER JOIN periods ON users.period_id = periods.id
INNER JOIN required_hours ON users.required_hour_id = required_hours.id
LEFT JOIN assistance_logs ON users.id = assistance_logs.user_id
WHERE users.is_admin = 0
  AND users.deleted_at IS NULL
GROUP BY 
    users.id, 
    users.name, 
    users.email, 
    periods.name, 
    required_hours.type, 
    required_hours.total_minutes
ORDER BY users.created_at ASC;