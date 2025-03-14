-- name: GenerateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
)
RETURNING *;

-- name: FindRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: FindUserByRefreshToken :one
SELECT * FROM users
WHERE id = (
    SELECT user_id FROM refresh_tokens
    WHERE token = $1 AND expires_at > $2 AND revoked_at IS NULL
);

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1 AND revoked_at IS NULL
RETURNING token, revoked_at;

