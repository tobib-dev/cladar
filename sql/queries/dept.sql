-- name: CreateDept :one
INSERT INTO departments (id, dept_name)
VALUES (gen_random_uuid(), $1)
RETURNING *;


-- name: GetAllDept :many
SELECT * FROM departments;

-- name: GetDeptByID :one
SELECT * FROM departments
WHERE dept_name = $1;