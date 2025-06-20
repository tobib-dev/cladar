// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: customers.sql

package database

import (
	"context"
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
