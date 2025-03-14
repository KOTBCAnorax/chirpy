-- name: UpgradeToRed :one
UPDATE users
SET is_chirpy_red = TRUE, updated_at = NOW()
WHERE id = $1
RETURNING id, is_chirpy_red, updated_at;