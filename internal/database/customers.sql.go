// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: customers.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createCustomer = `-- name: CreateCustomer :one
INSERT INTO customers (id, first_name, last_name, created_at, updated_at, email, phone, home, policy_type)
VALUES (gen_random_uuid(), $1, $2, NOW(), NOW(), $3, $4, $5, $6)
RETURNING id, first_name, last_name, created_at, updated_at, email, phone, home, policy_type
`

type CreateCustomerParams struct {
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	Home       string
	PolicyType string
}

func (q *Queries) CreateCustomer(ctx context.Context, arg CreateCustomerParams) (Customer, error) {
	row := q.db.QueryRowContext(ctx, createCustomer,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Phone,
		arg.Home,
		arg.PolicyType,
	)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Phone,
		&i.Home,
		&i.PolicyType,
	)
	return i, err
}

const deleteCustomer = `-- name: DeleteCustomer :exec
DELETE FROM customers
WHERE id = $1
`

func (q *Queries) DeleteCustomer(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteCustomer, id)
	return err
}

const getAllCustomers = `-- name: GetAllCustomers :many
SELECT id, first_name, last_name, created_at, updated_at, email, phone, home, policy_type FROM customers
`

func (q *Queries) GetAllCustomers(ctx context.Context) ([]Customer, error) {
	rows, err := q.db.QueryContext(ctx, getAllCustomers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Customer
	for rows.Next() {
		var i Customer
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Email,
			&i.Phone,
			&i.Home,
			&i.PolicyType,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCustomerByID = `-- name: GetCustomerByID :one
SELECT id, first_name, last_name, created_at, updated_at, email, phone, home, policy_type FROM customers
WHERE id=$1
`

func (q *Queries) GetCustomerByID(ctx context.Context, id uuid.UUID) (Customer, error) {
	row := q.db.QueryRowContext(ctx, getCustomerByID, id)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Phone,
		&i.Home,
		&i.PolicyType,
	)
	return i, err
}

const updateCustomer = `-- name: UpdateCustomer :one
UPDATE customers
SET first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    updated_at = NOW(),
    email = COALESCE($4, email),
    phone = COALESCE($5, phone),
    home = COALESCE($6, home),
    policy_type = COALESCE($7, policy_type)
WHERE id = $1
RETURNING id, first_name, last_name, created_at, updated_at, email, phone, home, policy_type
`

type UpdateCustomerParams struct {
	ID         uuid.UUID
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	Home       string
	PolicyType string
}

func (q *Queries) UpdateCustomer(ctx context.Context, arg UpdateCustomerParams) (Customer, error) {
	row := q.db.QueryRowContext(ctx, updateCustomer,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Phone,
		arg.Home,
		arg.PolicyType,
	)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Phone,
		&i.Home,
		&i.PolicyType,
	)
	return i, err
}
