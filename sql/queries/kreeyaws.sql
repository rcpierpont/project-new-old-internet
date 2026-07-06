-- name: CreateKreeyaw :one
INSERT INTO kreeyaws (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllKreeyaws :exec
DELETE FROM kreeyaws;

-- name: GetKreeyawsCreatedAtAsc :many
SELECT * FROM kreeyaws
ORDER BY created_at;

-- name: GetKreeyawsByAuthor :many
SELECT * FROM kreeyaws
WHERE user_id = $1
ORDER BY created_at;

-- name: GetKreeyaw :one
SELECT * FROM kreeyaws
WHERE id = $1;

-- name: DeleteKreeyaw :exec
DELETE FROM kreeyaws
WHERE id = $1;