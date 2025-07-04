-- name: CreateAgent :one
INSERT INTO agents (id, first_name, last_name, created_at, updated_at, email, dept)
VALUES (gen_random_uuid(), $1, $2, NOW(), NOW(), $3, $4)
RETURNING *;

-- name: GetAllAgents :many
SELECT * FROM agents;

-- name: GetAgentByID :one
SELECT * FROM agents
WHERE id = $1;

-- name: UpdateAgent :one
UPDATE agents
SET first_name = $2,
    last_name = $3,
    updated_at = NOW(),
    email = $4,
    dept = $5
WHERE id = $1
RETURNING *;
