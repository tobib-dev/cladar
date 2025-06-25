-- name: CreateCustomer :one
INSERT INTO customers (id, first_name, last_name, created_at, updated_at, email, phone, home, policy_type)
VALUES (gen_random_uuid(), $1, $2, NOW(), NOW(), $3, $4, $5, $6)
RETURNING *;

-- name: GetAllCustomers :many
SELECT * FROM customers;

-- name: GetCustomerByID :one
SELECT * FROM customers
WHERE id=$1;
