-- name: ListChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: ListUserChirps :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE id = $1;