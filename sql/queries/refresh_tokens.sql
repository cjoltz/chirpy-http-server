-- name: CreateToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1, -- token
    NOW(), -- created_at
    NOW(), -- updated_at
    $2, -- user_id
    $3, -- expires_at
    NULL --revoked_at
)
RETURNING *;

-- name: GetRefreshTokenByUserID :one
SELECT * FROM refresh_tokens
WHERE user_id = $1;

-- name: GetUserFromRefreshToken :one
SELECT * FROM refresh_tokens
WHERE 
    token = $1 
    AND revoked_at is NULL 
    AND expires_at > NOW();

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = NOW(), revoked_at = NOW()
WHERE token = $1;