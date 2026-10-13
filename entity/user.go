package entity

import "time"

type UserStatus string

const (
	UserStatusPending   UserStatus = "pending"
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusBlocked   UserStatus = "blocked"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

type UserRole string

const (
	UserRoleAdmin   UserRole = "admin"
	UserRoleUser    UserRole = "user"
)

// Entity
type User struct {
	ID             int64      `db:"id"`
	Name           string     `db:"name"`
	ProfileImg     string     `db:"profile_img"`
	Email          *string    `db:"email"`
	Phone          string     `db:"phone"`
	HashedPassword string     `db:"hashed_password"`
	Status         UserStatus `db:"status"`
	Address        string     `db:"address"`
	City           string     `db:"city"`
	CountryID      int64      `db:"country_id"`
	Country        *Country   `db:"-"`
	PostCode       string     `db:"post_code"`
	Role           UserRole   `db:"role"`
	FCMToken       *string    `db:"fcm_token"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
}
