-- +goose Up
CREATE TABLE departments (
    id UUID PRIMARY KEY,
    dept_name TEXT NOT NULL UNIQUE,
    manager_id UUID NOT NULL REFERENCES managers(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE departments;