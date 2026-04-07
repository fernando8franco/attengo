-- name: CreateRequiredHour :one
INSERT INTO required_hours (type, total_minutes) 
VALUES (?, ?)
RETURNING 
id,
type,
total_minutes;

-- name: DeleteRequiredHours :exec
DELETE FROM required_hours;

-- name: GetRequiredHours :many
SELECT id, type FROM required_hours;