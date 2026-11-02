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

type DesignationStatus string

const (
	DesignationStatusApproved DesignationStatus = "approved"
	DesignationStatusPending  DesignationStatus = "pending"
	DesignationStatusRejected DesignationStatus = "rejected"
)

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

// Entity
type User struct {
	ID                 int64              `db:"id"`
	Name               string             `db:"name"`
	ProfileImg         *string            `db:"profile_img"`
	Email              *string            `db:"email"`
	Phone              *string            `db:"phone"`
	HashedPassword     *string            `db:"hashed_password"`
	Status             UserStatus         `db:"status"`
	Address            *string            `db:"address"`
	City               *string            `db:"city"`
	CountryID          *int64             `db:"country_id"`
	CountryCode        *string            `db:"country_code"`
	PostCode           *string            `db:"post_code"`
	Role               UserRole           `db:"role"`
	DesignationID      *int64             `db:"designation_id"`
	DesignationTitle   *string            `db:"designation_title"`
	DesignationStatus  *DesignationStatus `db:"designation_status"`
	FirebaseUID        *string            `db:"firebase_uid"`
	PhoneVerifiedAt    *time.Time         `db:"phone_verified_at"`
	EmailVerifiedAt    *time.Time         `db:"email_verified_at"`
	ProfileCompletedAt *time.Time         `db:"profile_completed_at"`
	LastLoginAt        *time.Time         `db:"last_login_at"`
	CreatedAt          time.Time          `db:"created_at"`
	UpdatedAt          time.Time          `db:"updated_at"`
	DeletedAt          *time.Time         `db:"deleted_at"`
}
