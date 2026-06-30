-- name: CreateHonk :one
INSERT INTO honks (id, created_at, updated_at, body, post_id, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;