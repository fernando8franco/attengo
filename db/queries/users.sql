-- name: CreateUser :one
INSERT INTO users (id, name, email, password, required_hour_id, period_id) 
VALUES (?, ?, ?, ?, ?, ?)
RETURNING 
id,
name,
email,
password,
(SELECT type AS required_hour_type FROM required_hours WHERE id = users.required_hour_id);

-- name: ValidateUserPassword :one
SELECT COUNT(1) > 0
FROM users
WHERE is_admin = 0 AND id = ? AND password = ?;

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