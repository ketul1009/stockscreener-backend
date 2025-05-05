-- +goose up

ALTER TABLE screeners
ADD COLUMN rules JSONB NOT NULL,
ADD COLUMN username VARCHAR(255) NOT NULL REFERENCES users(username);

-- +goose Down

ALTER TABLE screeners
DROP COLUMN rules,
DROP COLUMN username;
