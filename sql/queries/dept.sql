-- name: CreateDept :one
INSERT INTO departments (id, dept_name, created_at, updated_at)
VALUES (gen_random_uuid(), $1, NOW(), NOW())
RETURNING *;


-- name: GetAllDept :many
SELECT * FROM departments;

-- name: GetDeptByID :one
SELECT * FROM departments
WHERE dept_name = $1;
