package entity

import "time"

type Permission struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Module      string `db:"module"`
	Action      string `db:"action"`
}

type UserRole string

const (
	UserRoleSuperAdmin UserRole = "super_admin"
	UserRoleStaff      UserRole = "staff"
	UserRoleAdmin      UserRole = "admin"
	UserRoleUser       UserRole = "user"
)

type Scope string

const (
	ScopeAll      Scope = "all"
	ScopeAdmin    Scope = "admin"
	ScopePlatform Scope = "platform"
)

// Entity
type Role struct {
	ID          int64        `db:"id"`
	Name        UserRole     `db:"name"`
	Scope       Scope        `db:"scope"`
	Permissions []Permission `db:"-"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
}
