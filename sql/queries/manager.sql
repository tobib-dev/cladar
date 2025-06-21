-- name: CreateManager :one
INSERT INTO managers (id, first_name, last_name, email, dept_id)
VALUES (gen_random_uuid(), $1, $2, $3, $4)
RETURNING *;
