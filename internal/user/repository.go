package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Sabbir185/gpi/entity"
)

type NewUser struct {
	Name            string
	Email           *string
	Phone           *string
	ProfileImg      *string
	FirebaseUID     *string
	PhoneVerifiedAt *time.Time
	EmailVerifiedAt *time.Time
}

type ProfileUpdate struct {
	Name              string
	CountryCode       string
	ProfileImg        *string
	DesignationID     *int64
	DesignationTitle  *string
	DesignationStatus *entity.DesignationStatus
}

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	FindByPhone(ctx context.Context, phone string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByFirebaseUID(ctx context.Context, uid string) (*entity.User, error)
	Create(ctx context.Context, user *NewUser) (*entity.User, error)
	LinkFirebase(ctx context.Context, id int64, uid string, emailVerified bool) error
	MarkPhoneVerified(ctx context.Context, id int64) error
	TouchLastLogin(ctx context.Context, id int64) error
	UpdateProfile(ctx context.Context, id int64, profile *ProfileUpdate) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

const selectUser = `
	SELECT u.id, u.name, u.profile_img, u.email, u.phone, u.status, u.country_id, u.country_code, u.role,
		u.designation_id, u.designation_title, u.designation_status,
		u.firebase_uid, u.phone_verified_at, u.email_verified_at, u.profile_completed_at,
		u.last_login_at, u.created_at, u.updated_at
	FROM users u
`

func (r *userRepository) findOne(ctx context.Context, where string, arg any) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRowContext(ctx, selectUser+where+" AND u.deleted_at IS NULL LIMIT 1;", arg).Scan(
		&user.ID,
		&user.Name,
		&user.ProfileImg,
		&user.Email,
		&user.Phone,
		&user.Status,
		&user.CountryID,
		&user.CountryCode,
		&user.Role,
		&user.DesignationID,
		&user.DesignationTitle,
		&user.DesignationStatus,
		&user.FirebaseUID,
		&user.PhoneVerifiedAt,
		&user.EmailVerifiedAt,
		&user.ProfileCompletedAt,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	return r.findOne(ctx, "WHERE u.id = $1", id)
}

func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	return r.findOne(ctx, "WHERE u.phone = $1", phone)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return r.findOne(ctx, "WHERE LOWER(u.email) = LOWER($1)", email)
}

func (r *userRepository) FindByFirebaseUID(ctx context.Context, uid string) (*entity.User, error) {
	return r.findOne(ctx, "WHERE u.firebase_uid = $1", uid)
}

// Create returns nil, nil when a concurrent request created the same account first.
func (r *userRepository) Create(ctx context.Context, user *NewUser) (*entity.User, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users(name, email, phone, profile_img, firebase_uid, phone_verified_at, email_verified_at, status, role)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT DO NOTHING
		RETURNING id;
	`,
		user.Name,
		user.Email,
		user.Phone,
		user.ProfileImg,
		user.FirebaseUID,
		user.PhoneVerifiedAt,
		user.EmailVerifiedAt,
		entity.UserStatusActive,
		entity.UserRoleUser,
	).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *userRepository) LinkFirebase(ctx context.Context, id int64, uid string, emailVerified bool) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET
			firebase_uid = $2,
			email_verified_at = CASE WHEN $3 THEN COALESCE(email_verified_at, NOW()) ELSE email_verified_at END,
			updated_at = NOW()
		WHERE id = $1;
	`, id, uid, emailVerified)
	return err
}

func (r *userRepository) MarkPhoneVerified(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET phone_verified_at = COALESCE(phone_verified_at, NOW()), updated_at = NOW()
		WHERE id = $1;
	`, id)
	return err
}

func (r *userRepository) TouchLastLogin(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET last_login_at = NOW() WHERE id = $1;`, id)
	return err
}

func (r *userRepository) UpdateProfile(ctx context.Context, id int64, profile *ProfileUpdate) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET
			name = $2,
			country_code = $3,
			profile_img = $4,
			designation_id = $5,
			designation_title = $6,
			designation_status = $7,
			profile_completed_at = COALESCE(profile_completed_at, NOW()),
			updated_at = NOW()
		WHERE id = $1;
	`,
		id,
		profile.Name,
		profile.CountryCode,
		profile.ProfileImg,
		profile.DesignationID,
		profile.DesignationTitle,
		profile.DesignationStatus,
	)
	return err
}
