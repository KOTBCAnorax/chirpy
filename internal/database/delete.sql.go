// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: delete.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const deleteChirp = `-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1 AND user_id = $2
`

type DeleteChirpParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteChirp(ctx context.Context, arg DeleteChirpParams) error {
	_, err := q.db.ExecContext(ctx, deleteChirp, arg.ID, arg.UserID)
	return err
}

const isOwner = `-- name: IsOwner :one
SELECT EXISTS(
    SELECT id, created_at, updated_at, body, user_id FROM chirps
    WHERE id = $1 AND user_id = $2
)AS is_owner
`

type IsOwnerParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) IsOwner(ctx context.Context, arg IsOwnerParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, isOwner, arg.ID, arg.UserID)
	var is_owner bool
	err := row.Scan(&is_owner)
	return is_owner, err
}
