package main

import "github.com/google/uuid"

type UserType string

const (
	UserRoleAgent   UserType = "agent"
	UserRoleManager UserType = "manager"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	UserRole UserType  `json:"user_role"`
	RoleID   uuid.UUID `json:"role_id"`
}
