package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Sabbir185/gpi/entity"
	"github.com/google/uuid"
)

type AuthRepository interface {
	FindLatestOTP(ctx context.Context, identity string, otpType entity.OtpType) (*entity.Otp, error)
	CountOTPsSince(ctx context.Context, identity string, otpType entity.OtpType, since time.Time) (int, error)
	CreateOTP(ctx context.Context, identity, codeHash string, otpType entity.OtpType, channel entity.OtpChannel, ttl time.Duration) (*entity.Otp, error)
	ExpireOTP(ctx context.Context, id int64) error
	ConsumeOTPAttempt(ctx context.Context, requestID, identity string, otpType entity.OtpType, maxAttempts int) (int64, string, error)
	MarkOTPConsumed(ctx context.Context, id int64) (bool, error)

	CreateSession(ctx context.Context, session *entity.Session) error
	FindSessionByRefreshHash(ctx context.Context, hash string) (*entity.Session, error)
	FindSessionByPreviousRefreshHash(ctx context.Context, hash string) (*entity.Session, error)
	RotateSession(ctx context.Context, id, currentHash, newHash string, expiresAt time.Time) (bool, error)
	RevokeSession(ctx context.Context, id string) error
	FindSession(ctx context.Context, id string, userID int64) (*entity.Session, error)
	ValidateSession(ctx context.Context, sessionID string, userID int64) (bool, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindLatestOTP(ctx context.Context, identity string, otpType entity.OtpType) (*entity.Otp, error) {
	var otp entity.Otp
	err := r.db.QueryRowContext(ctx, `
		SELECT id, request_id, identity, code_hash, type, channel, attempt, expires_at, consumed_at, created_at, updated_at
		FROM otps
		WHERE identity = $1 AND type = $2
		ORDER BY created_at DESC
		LIMIT 1;
	`, identity, otpType).Scan(
		&otp.ID,
		&otp.RequestID,
		&otp.Identity,
		&otp.CodeHash,
		&otp.Type,
		&otp.Channel,
		&otp.Attempt,
		&otp.ExpiresAt,
		&otp.ConsumedAt,
		&otp.CreatedAt,
		&otp.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &otp, nil
}

func (r *authRepository) CountOTPsSince(ctx context.Context, identity string, otpType entity.OtpType, since time.Time) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM otps WHERE identity = $1 AND type = $2 AND created_at > $3;
	`, identity, otpType, since).Scan(&count)
	return count, err
}

// CreateOTP expires any code still pending for the identity, so only the newest one works.
func (r *authRepository) CreateOTP(ctx context.Context, identity, codeHash string, otpType entity.OtpType, channel entity.OtpChannel, ttl time.Duration) (*entity.Otp, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE otps SET expires_at = NOW(), updated_at = NOW()
		WHERE identity = $1 AND type = $2 AND expires_at > NOW() AND consumed_at IS NULL;
	`, identity, otpType)
	if err != nil {
		return nil, err
	}

	otp := entity.Otp{RequestID: uuid.NewString(), Identity: identity, CodeHash: codeHash, Type: otpType, Channel: channel}
	err = tx.QueryRowContext(ctx, `
		INSERT INTO otps(request_id, identity, code_hash, type, channel, expires_at)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id, attempt, expires_at, created_at, updated_at;
	`, otp.RequestID, identity, codeHash, otpType, channel, time.Now().UTC().Add(ttl)).Scan(
		&otp.ID,
		&otp.Attempt,
		&otp.ExpiresAt,
		&otp.CreatedAt,
		&otp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &otp, tx.Commit()
}

func (r *authRepository) ExpireOTP(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE otps SET expires_at = NOW(), updated_at = NOW() WHERE id = $1;`, id)
	return err
}

// ConsumeOTPAttempt counts the attempt before the code is compared, so parallel
// guesses can't exceed maxAttempts. It returns an empty hash when the code is expired, used or used up.
func (r *authRepository) ConsumeOTPAttempt(ctx context.Context, requestID, identity string, otpType entity.OtpType, maxAttempts int) (int64, string, error) {
	var id int64
	var codeHash string
	err := r.db.QueryRowContext(ctx, `
		UPDATE otps SET attempt = attempt + 1, updated_at = NOW()
		WHERE request_id = $1 AND identity = $2 AND type = $3
		  AND expires_at > NOW() AND consumed_at IS NULL AND attempt < $4
		RETURNING id, code_hash;
	`, requestID, identity, otpType, maxAttempts).Scan(&id, &codeHash)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, "", nil
	}
	return id, codeHash, err
}

// MarkOTPConsumed keeps the row so it still counts towards the hourly send limit.
func (r *authRepository) MarkOTPConsumed(ctx context.Context, id int64) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE otps SET consumed_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND consumed_at IS NULL;
	`, id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	return affected == 1, err
}

func (r *authRepository) CreateSession(ctx context.Context, session *entity.Session) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO sessions(id, user_id, device_id, refresh_token_hash, expires_at)
		VALUES($1, $2, $3, $4, $5);
	`, session.ID, session.UserID, session.DeviceID, session.RefreshTokenHash, session.ExpiresAt)
	return err
}

const selectSession = `
	SELECT id, user_id, device_id, refresh_token_hash, previous_refresh_token_hash, rotated_at,
		expires_at, revoked_at, last_used_at, created_at, updated_at
	FROM sessions
`

func (r *authRepository) findSession(ctx context.Context, where string, args ...any) (*entity.Session, error) {
	var s entity.Session
	err := r.db.QueryRowContext(ctx, selectSession+where, args...).Scan(
		&s.ID,
		&s.UserID,
		&s.DeviceID,
		&s.RefreshTokenHash,
		&s.PreviousRefreshTokenHash,
		&s.RotatedAt,
		&s.ExpiresAt,
		&s.RevokedAt,
		&s.LastUsedAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *authRepository) FindSessionByRefreshHash(ctx context.Context, hash string) (*entity.Session, error) {
	return r.findSession(ctx, "WHERE refresh_token_hash = $1;", hash)
}

func (r *authRepository) FindSessionByPreviousRefreshHash(ctx context.Context, hash string) (*entity.Session, error) {
	return r.findSession(ctx, "WHERE previous_refresh_token_hash = $1;", hash)
}

func (r *authRepository) FindSession(ctx context.Context, id string, userID int64) (*entity.Session, error) {
	return r.findSession(ctx, "WHERE id = $1 AND user_id = $2;", id, userID)
}

// RotateSession only succeeds for the caller still holding currentHash, which makes a
// rotation race lose cleanly instead of issuing two valid refresh tokens.
func (r *authRepository) RotateSession(ctx context.Context, id, currentHash, newHash string, expiresAt time.Time) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE sessions SET
			previous_refresh_token_hash = CASE WHEN refresh_token_hash = $2 THEN $2 ELSE previous_refresh_token_hash END,
			refresh_token_hash = $3,
			rotated_at = CASE WHEN refresh_token_hash = $2 THEN NOW() ELSE rotated_at END,
			expires_at = $4,
			last_used_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND revoked_at IS NULL AND (refresh_token_hash = $2 OR previous_refresh_token_hash = $2);
	`, id, currentHash, newHash, expiresAt)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	return affected == 1, err
}

func (r *authRepository) RevokeSession(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE sessions SET revoked_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND revoked_at IS NULL;
	`, id)
	return err
}

func (r *authRepository) ValidateSession(ctx context.Context, sessionID string, userID int64) (bool, error) {
	var valid bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM sessions s
			JOIN users u ON u.id = s.user_id
			WHERE s.id = $1 AND s.user_id = $2
			  AND s.revoked_at IS NULL AND s.expires_at > NOW()
			  AND u.status = $3 AND u.deleted_at IS NULL
		);
	`, sessionID, userID, entity.UserStatusActive).Scan(&valid)
	return valid, err
}
