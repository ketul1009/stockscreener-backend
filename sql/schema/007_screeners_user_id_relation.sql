-- +goose up
ALTER TABLE screeners
DROP CONSTRAINT unique_name_user_id,
DROP COLUMN username,
ADD COLUMN user_id UUID NOT NULL REFERENCES users(id),
ADD CONSTRAINT unique_name_user_id UNIQUE (name, user_id);

-- +goose Down
ALTER TABLE screeners
DROP CONSTRAINT unique_name_user_id,
DROP COLUMN user_id,
ADD COLUMN username VARCHAR(255) NOT NULL REFERENCES users(username),
ADD CONSTRAINT unique_name_user_id UNIQUE (name, username);
