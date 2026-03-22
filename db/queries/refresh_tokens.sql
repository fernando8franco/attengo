-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at)
VALUES (?, ?, datetime('now', '+30 days'))
RETURNING *;

-- name: GetUserIdFromRefreshToken :one
SELECT user_id FROM refresh_tokens
WHERE token = ?
AND is_revoked = 0
AND deleted_at IS NULL
AND datetime(expires_at) > datetime('now');

-- name: SetRevokedAt :exec
UPDATE refresh_tokens
SET is_revoked = 1, updated_at = datetime('now')
WHERE token = ?;

-- name: GetRevoked :one
SELECT is_revoked FROM refresh_tokens;