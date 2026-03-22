-- name: CreateUser :one
INSERT INTO users (id, name, email, password, required_hour_id, period_id) 
VALUES (?, ?, ?, ?, ?, ?)
RETURNING 
id;

-- name: GetUsersPasswords :many
SELECT password FROM users WHERE is_admin = 0;

-- name: ValidateUserPassword :one
SELECT EXISTS (
  SELECT 1
  FROM users
  WHERE is_admin = 0 AND id = ? AND password = ?
) = 1;

-- name: CreateAdmin :one
INSERT INTO users (id, is_admin, name, email, password) 
VALUES (?, 1, ?, ?, ?)
RETURNING 
id,
name,
email;

-- name: ExistsAdmin :one
SELECT EXISTS (
  SELECT 1
  FROM users
  WHERE is_admin = 1
) = 1;