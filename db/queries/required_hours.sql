-- name: CreateRequiredHours :one
INSERT INTO required_hours (type, minutes, created_at, updated_at) 
VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING *;