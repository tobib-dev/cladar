-- name: CreateUser :one
INSERT INTO users (id, email, pswd, user_role, role_id)
VALUES (gen_random_uuid(), $1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;