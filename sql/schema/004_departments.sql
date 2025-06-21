-- +goose Up
CREATE TABLE departments (
    id UUID PRIMARY KEY,
    dept_name TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE departments;