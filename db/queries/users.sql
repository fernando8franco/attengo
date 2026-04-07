-- name: CreateUser :one
INSERT INTO users (id, name, email, password, required_hour_id, period_id) 
VALUES (?, ?, ?, ?, ?, ?)
RETURNING 
id;

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