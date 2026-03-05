-- name: CreateRequiredHours :one
INSERT INTO required_hours (type, total_minutes, created_at, updated_at) 
VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING 
id,
type,
total_minutes;

-- name: DeleteRequiredHours :exec
DELETE FROM required_hours;