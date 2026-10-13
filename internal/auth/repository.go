package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/Sabbir185/gpi/entity"
)

type AuthRepository interface {
	CreateOTP(ctx context.Context, identity string, otp string, OtpType entity.OtpType, OtpChannel entity.OtpChannel, ttl time.Duration) (*entity.Otp, error)
	FindOTPByIdentityAndType(identity string, OtpType entity.OtpType) (*entity.Otp, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

func (s *authRepository) CreateOTP(ctx context.Context, identity string, otp string, OtpType entity.OtpType, OtpChannel entity.OtpChannel, ttl time.Duration) (*entity.Otp, error) {
	expireAt := time.Now().UTC().Add(ttl)
	query := `
		INSERT INTO otps(identity, otp, type, channel, expires_at)
		VALUES($1, $2, $3, $4, $5)
		RETURNING *;
	`
	row := s.db.QueryRowContext(ctx, query, identity, otp, OtpType, OtpChannel, expireAt)
	if err := row.Err(); err != nil {
		return nil, err
	}
	var result entity.Otp
	err := row.Scan(
		&result.ID,
		&result.Identity,
		&result.OTP,
		&result.Type,
		&result.Channel,
		&result.Attempt,
		&result.ExpiresAt,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *authRepository) FindOTPByIdentityAndType(identity string, OtpType entity.OtpType) (*entity.Otp, error) {
	return nil, nil
}
