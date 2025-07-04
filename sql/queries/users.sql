-- name: CreateUser :one
INSERT INTO users (id, email, pswd, user_role, role_id, created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, $3, $4, NOW(), NOW())
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: DeleteUserByRoleId :exec
DELETE FROM users
WHERE role_id = $1;
