-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, user_id, body) VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2
)
RETURNING *;

-- name: DeleteChirps :exec
DELETE FROM chirps;

-- name: GetAllChirps :many
SELECT * FROM chirps ORDER BY created_at ASC;

-- name: GetChirpsForUser :many
SELECT * FROM chirps WHERE user_id = $1 ORDER BY created_at ASC;

-- name: GetChirpByID :one
SELECT * FROM chirps WHERE id = $1;

-- name: GetChirpForUserByID :one
SELECT count(*) FROM chirps WHERE id = $1 and user_id = $2;

-- name: DeleteChirpForUserByID :exec
DELETE FROM chirps WHERE id = $1 and user_id = $2;
