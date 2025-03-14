-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1 AND user_id = $2;

-- name: IsOwner :one
SELECT EXISTS(
    SELECT * FROM chirps
    WHERE id = $1 AND user_id = $2
)AS is_owner;