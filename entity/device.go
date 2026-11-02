package entity

import "time"

type Device struct {
	ID             int64     `db:"id"`
	UserID         int64     `db:"user_id"`
	InstallationID string    `db:"installation_id"`
	Platform       string    `db:"platform"`
	OSVersion      *string   `db:"os_version"`
	AppVersion     *string   `db:"app_version"`
	Locale         *string   `db:"locale"`
	Timezone       *string   `db:"timezone"`
	PushToken      *string   `db:"push_token"`
	LastSeenAt     time.Time `db:"last_seen_at"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
