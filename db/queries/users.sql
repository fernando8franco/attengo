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