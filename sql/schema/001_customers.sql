-- +goose Up
CREATE TABLE customers (
    id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    home TEXT NOT NULL,
    policy_type TEXT NOT NULL
);

-- +goose Down
DROP TABLE customers;