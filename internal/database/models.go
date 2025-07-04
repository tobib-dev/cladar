// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusDeclined  Status = "declined"
	StatusAwarded   Status = "awarded"
	StatusCompleted Status = "completed"
	StatusPending   Status = "pending"
)

func (e *Status) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Status(s)
	case string:
		*e = Status(s)
	default:
		return fmt.Errorf("unsupported scan type for Status: %T", src)
	}
	return nil
}

type NullStatus struct {
	Status Status
	Valid  bool // Valid is true if Status is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullStatus) Scan(value interface{}) error {
	if value == nil {
		ns.Status, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Status.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Status), nil
}

type UserType string

const (
	UserTypeAgent   UserType = "agent"
	UserTypeManager UserType = "manager"
)

func (e *UserType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserType(s)
	case string:
		*e = UserType(s)
	default:
		return fmt.Errorf("unsupported scan type for UserType: %T", src)
	}
	return nil
}

type NullUserType struct {
	UserType UserType
	Valid    bool // Valid is true if UserType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUserType) Scan(value interface{}) error {
	if value == nil {
		ns.UserType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UserType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUserType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UserType), nil
}

type Agent struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
	Dept      string
}

type Claim struct {
	ID            uuid.UUID
	CustomerID    uuid.UUID
	AgentID       uuid.UUID
	ClaimType     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CurrentStatus Status
	Award         sql.NullFloat64
}

type Customer struct {
	ID         uuid.UUID
	FirstName  string
	LastName   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Email      string
	Phone      string
	Home       string
	PolicyType string
}

type Department struct {
	ID        uuid.UUID
	DeptName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Manager struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	DeptID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RefreshToken struct {
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt sql.NullTime
}

type User struct {
	ID        uuid.UUID
	Email     string
	Pswd      string
	UserRole  UserType
	RoleID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
