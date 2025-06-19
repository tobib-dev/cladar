-- name: CreateAgent :one
INSERT INTO agents (id, first_name, last_name, created_at, updated_at, dept)
VALUES (gen_random_uuid(), $1, $2, NOW(), NOW(), $3)
RETURNING *;