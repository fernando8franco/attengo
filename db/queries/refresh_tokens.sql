-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token, user_id, expires_at)
VALUES (?, ?, datetime('now', '+30 days'));

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