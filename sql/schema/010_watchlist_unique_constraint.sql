-- +goose Up
ALTER TABLE watchlist ADD CONSTRAINT unique_watchlist_name UNIQUE (name, user_id);

-- +goose Down
ALTER TABLE watchlist DROP CONSTRAINT unique_watchlist_name;

