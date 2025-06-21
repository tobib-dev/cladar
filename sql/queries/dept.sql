-- name: CreateDept :one
INSERT INTO departments (id, dept_name)
VALUES (gen_random_uuid(), $1)
RETURNING *;
