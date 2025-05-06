-- +goose Up
ALTER TABLE screeners ADD COLUMN stock_universe TEXT NOT NULL DEFAULT 'nifty50';
ALTER TABLE screeners ALTER COLUMN name SET NOT NULL;

-- +goose Down
ALTER TABLE screeners DROP COLUMN stock_universe;
ALTER TABLE screeners ALTER COLUMN name DROP NOT NULL;

