-- name: StoreRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES ($1, NOW(), NOW(), $2, $3, $4)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: UpdateRefreshToken :one
UPDATE refresh_tokens
SET token = $2,
    updated_at = NOW()
WHERE token = $1
RETURNING *;

-- name: RevokeToken :one
UPDATE refresh_tokens
SET revoked_at = NOW(),
    updated_at = NOW()
where token = $1
RETURNING *;
