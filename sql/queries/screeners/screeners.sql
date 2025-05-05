-- name: CreateScreener :one
INSERT INTO screeners (username, name, rules)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetScreener :one
SELECT * FROM screeners
WHERE id = $1;




