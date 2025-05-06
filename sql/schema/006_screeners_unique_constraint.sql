-- +goose Up
ALTER TABLE screeners ADD CONSTRAINT unique_name_user_id UNIQUE (name, username);

-- +goose Down
ALTER TABLE screeners DROP CONSTRAINT unique_name_user_id;


