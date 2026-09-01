package entity

import (
	"time"

	"github.com/Sabbir185/geopunch/infra/db"
)

type Permission struct {
	ID          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Module      string `db:"module" json:"module"`
	Action      string `db:"action" json:"action"`
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
	ID          int64                   `db:"id"`
	Name        UserRole                `db:"name"`
	Scope       Scope                   `db:"scope"`
	Permissions db.JSONB[[]Permission] `db:"permissions"`
	CreatedAt   time.Time               `db:"created_at"`
	UpdatedAt   time.Time               `db:"updated_at"`
}
