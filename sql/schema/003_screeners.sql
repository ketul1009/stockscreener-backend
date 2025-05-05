-- +goose up

CREATE TABLE screeners (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- +goose Down

DROP TABLE screeners;

