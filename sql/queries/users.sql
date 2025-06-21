-- name: CreateUser :one
INSERT INTO users (id, email, pswd, user_role, role_id)
VALUES (gen_random_uuid(), $1, $2, $3, $4)
RETURNING *;