-- name: GetChirp :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: GetChirpByUserID :many
SELECT *
FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;
