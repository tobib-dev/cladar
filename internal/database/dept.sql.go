// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: dept.sql

package database

import (
	"context"
)

const createDept = `-- name: CreateDept :one
INSERT INTO departments (id, dept_name, created_at, updated_at)
VALUES (gen_random_uuid(), $1, NOW(), NOW())
RETURNING id, dept_name, created_at, updated_at
`

func (q *Queries) CreateDept(ctx context.Context, deptName string) (Department, error) {
	row := q.db.QueryRowContext(ctx, createDept, deptName)
	var i Department
	err := row.Scan(
		&i.ID,
		&i.DeptName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllDept = `-- name: GetAllDept :many
SELECT id, dept_name, created_at, updated_at FROM departments
`

func (q *Queries) GetAllDept(ctx context.Context) ([]Department, error) {
	rows, err := q.db.QueryContext(ctx, getAllDept)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Department
	for rows.Next() {
		var i Department
		if err := rows.Scan(
			&i.ID,
			&i.DeptName,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getDeptByID = `-- name: GetDeptByID :one
SELECT id, dept_name, created_at, updated_at FROM departments
WHERE dept_name = $1
`

func (q *Queries) GetDeptByID(ctx context.Context, deptName string) (Department, error) {
	row := q.db.QueryRowContext(ctx, getDeptByID, deptName)
	var i Department
	err := row.Scan(
		&i.ID,
		&i.DeptName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
