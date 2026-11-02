package device

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrOwnedByAnotherUser = errors.New("device belongs to another user")

type PushTarget struct {
	Token  string
	Locale string
}

type ActiveDevice struct {
	ID           int64
	Platform     string
	OSVersion    *string
	LastActiveAt time.Time
}

type DeviceRepository interface {
	Upsert(ctx context.Context, userID int64, device *DeviceRequest, allowTransfer bool) (int64, error)
	SessionDeviceID(ctx context.Context, sessionID string) (*int64, error)
	ListActiveForUser(ctx context.Context, userID int64) ([]ActiveDevice, error)
	ListPushTargets(ctx context.Context, userID int64) ([]PushTarget, error)
	ClearPushToken(ctx context.Context, deviceID int64) error
	DeletePushToken(ctx context.Context, token string) error
	RevokeSessions(ctx context.Context, userID, deviceID int64) (int64, error)
}

type deviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) DeviceRepository {
	return &deviceRepository{db: db}
}

// Upsert attaches the install to the user and moves its push token here, so a
// token is never registered on two records and pushes can't reach a previous account.
// Only a fresh sign-in may move an install from one account to another (allowTransfer).
func (r *deviceRepository) Upsert(ctx context.Context, userID int64, device *DeviceRequest, allowTransfer bool) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if device.PushToken != nil {
		_, err = tx.ExecContext(ctx, `
			UPDATE devices SET push_token = NULL, updated_at = NOW()
			WHERE push_token = $1 AND installation_id <> $2;
		`, *device.PushToken, device.InstallationID)
		if err != nil {
			return 0, err
		}
	}

	var id int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO devices(user_id, installation_id, platform, os_version, app_version, locale, timezone, push_token)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (installation_id) DO UPDATE SET
			user_id = EXCLUDED.user_id,
			platform = EXCLUDED.platform,
			os_version = EXCLUDED.os_version,
			app_version = EXCLUDED.app_version,
			locale = EXCLUDED.locale,
			timezone = EXCLUDED.timezone,
			push_token = EXCLUDED.push_token,
			last_seen_at = NOW(),
			updated_at = NOW()
		WHERE $9 OR devices.user_id = EXCLUDED.user_id
		RETURNING id;
	`,
		userID,
		device.InstallationID,
		device.Platform,
		nullIfEmpty(device.OSVersion),
		device.AppVersion,
		nullIfEmpty(device.Locale),
		device.Timezone,
		device.PushToken,
		allowTransfer,
	).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrOwnedByAnotherUser
	}
	if err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

func (r *deviceRepository) SessionDeviceID(ctx context.Context, sessionID string) (*int64, error) {
	var deviceID *int64
	err := r.db.QueryRowContext(ctx, `SELECT device_id FROM sessions WHERE id = $1;`, sessionID).Scan(&deviceID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return deviceID, err
}

func (r *deviceRepository) ListActiveForUser(ctx context.Context, userID int64) ([]ActiveDevice, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT d.id, d.platform, d.os_version, GREATEST(d.last_seen_at, MAX(s.last_used_at))
		FROM devices d
		JOIN sessions s ON s.device_id = d.id
		WHERE s.user_id = $1 AND s.revoked_at IS NULL AND s.expires_at > NOW()
		GROUP BY d.id
		ORDER BY 4 DESC;
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []ActiveDevice
	for rows.Next() {
		var d ActiveDevice
		if err := rows.Scan(&d.ID, &d.Platform, &d.OSVersion, &d.LastActiveAt); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, rows.Err()
}

func (r *deviceRepository) ListPushTargets(ctx context.Context, userID int64) ([]PushTarget, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT push_token, COALESCE(locale, '')
		FROM devices
		WHERE user_id = $1 AND push_token IS NOT NULL;
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []PushTarget
	for rows.Next() {
		var t PushTarget
		if err := rows.Scan(&t.Token, &t.Locale); err != nil {
			return nil, err
		}
		targets = append(targets, t)
	}
	return targets, rows.Err()
}

func (r *deviceRepository) ClearPushToken(ctx context.Context, deviceID int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE devices SET push_token = NULL, updated_at = NOW() WHERE id = $1;`, deviceID)
	return err
}

func (r *deviceRepository) DeletePushToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE devices SET push_token = NULL, updated_at = NOW() WHERE push_token = $1;`, token)
	return err
}

func (r *deviceRepository) RevokeSessions(ctx context.Context, userID, deviceID int64) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE sessions SET revoked_at = NOW(), updated_at = NOW()
		WHERE user_id = $1 AND device_id = $2 AND revoked_at IS NULL;
	`, userID, deviceID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func nullIfEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
