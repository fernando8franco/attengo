-- name: CreateUser :one
INSERT INTO users (id, name, email, password, required_hour_id) 
VALUES (?, ?, ?, ?, ?)
RETURNING 
id,
name,
email,
password,
(SELECT type AS required_hour_type FROM required_hours WHERE id = users.required_hour_id),
(SELECT total_minutes AS required_hour_minutes FROM required_hours WHERE id = users.required_hour_id);

-- name: ValidateUserPassword :one
SELECT COUNT(1) > 0
FROM users
WHERE id = ? AND password = ?;