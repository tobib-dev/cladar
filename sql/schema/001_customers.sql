-- +goose Up
CREATE TABLE customers (
    id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone TEXT NOT NULL UNIQUE,
    home TEXT NOT NULL UNIQUE,
    policy_type TEXT NOT NULL
);

-- +goose Down
DROP TABLE customers;