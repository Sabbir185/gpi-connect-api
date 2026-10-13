package entity

import (
	"time"

	"github.com/Sabbir185/gpi/infra/db"
)

type Permission struct {
	ID          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Module      string `db:"module" json:"module"`
	Action      string `db:"action" json:"action"`
}

type RoleName string

const (
	RoleNameSuperAdmin RoleName = "super_admin"
	RoleNameStaff      RoleName = "staff"
	RoleNameAdmin      RoleName = "admin"
	RoleNameUser       RoleName = "user"
)

type Scope string

const (
	ScopeAll      Scope = "all"
	ScopeAdmin    Scope = "admin"
	ScopePlatform Scope = "platform"
)

// Entity
type Role struct {
	ID          int64                  `db:"id"`
	Name        RoleName               `db:"name"`
	Scope       Scope                  `db:"scope"`
	Permissions db.JSONB[[]Permission] `db:"permissions"`
	CreatedAt   time.Time              `db:"created_at"`
	UpdatedAt   time.Time              `db:"updated_at"`
}
