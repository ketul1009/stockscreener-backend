-- name: CreateScreener :one
INSERT INTO screeners (username, name, rules)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetScreener :one
SELECT * FROM screeners
WHERE id = $1;

-- name: GetScreeners :many
SELECT * FROM screeners
WHERE username = $1;


