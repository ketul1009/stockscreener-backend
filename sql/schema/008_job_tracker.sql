-- +goose Up

CREATE TABLE job_tracker (
    id SERIAL PRIMARY KEY,
    job_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    job_status VARCHAR(255) NOT NULL,
    job_created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    job_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down

DROP TABLE job_tracker;

