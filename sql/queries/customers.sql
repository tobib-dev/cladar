-- name: CreateCustomer :one
INSERT INTO customers (id, first_name, last_name, created_at, updated_at, email, phone, home, policy_type)
VALUES (gen_random_uuid(), $1, $2, NOW(), NOW(), $3, $4, $5, $6)
RETURNING *;

-- name: GetAllCustomers :many
SELECT * FROM customers;

-- name: GetCustomerByID :one
SELECT * FROM customers
WHERE id=$1;

-- name: UpdateCustomer :one
UPDATE customers
SET first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    updated_at = NOW(),
    email = COALESCE($4, email),
    phone = COALESCE($5, phone),
    home = COALESCE($6, home),
    policy_type = COALESCE($7, policy_type)
WHERE id = $1
RETURNING *;

-- name: DeleteCustomer :exec
DELETE FROM customers
WHERE id = $1;
