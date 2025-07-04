-- name: CreateManager :one
INSERT INTO managers (id, first_name, last_name, email, dept_id, created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, $3, $4, NOW(), NOW())
RETURNING *;

-- name: GetAllManagers :many
SELECT * FROM managers;
