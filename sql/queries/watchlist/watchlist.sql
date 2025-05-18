-- name: CreateWatchlist :one

INSERT INTO watchlist (name, user_id, stocks, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: GetWatchlist :one

SELECT * FROM watchlist
WHERE id = $1;

-- name: GetAllWatchlists :many

SELECT * FROM watchlist
WHERE user_id = $1;

-- name: UpdateWatchlist :one

UPDATE watchlist
SET name = $3, stocks = $4, updated_at = $5
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteWatchlist :exec

DELETE FROM watchlist
WHERE id = $1;
