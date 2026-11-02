package entity

import "time"

type Session struct {
	ID                       string     `db:"id"`
	UserID                   int64      `db:"user_id"`
	DeviceID                 *int64     `db:"device_id"`
	RefreshTokenHash         string     `db:"refresh_token_hash"`
	PreviousRefreshTokenHash *string    `db:"previous_refresh_token_hash"`
	RotatedAt                *time.Time `db:"rotated_at"`
	ExpiresAt                time.Time  `db:"expires_at"`
	RevokedAt                *time.Time `db:"revoked_at"`
	LastUsedAt               time.Time  `db:"last_used_at"`
	CreatedAt                time.Time  `db:"created_at"`
	UpdatedAt                time.Time  `db:"updated_at"`
}
