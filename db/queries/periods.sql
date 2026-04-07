-- name: CreatePeriod :one
INSERT INTO periods (name, entry_date, exit_date)
VALUES (?, ?, ?)
RETURNING
id,
name,
entry_date,
exit_date;

-- name: DeletePeriods :exec
DELETE FROM periods;

-- name: GetPeriods :many
SELECT id, name FROM periods;