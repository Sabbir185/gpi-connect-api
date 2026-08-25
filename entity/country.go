package entity

import "time"

type CountryStatus string

const (
	CountryStatusActive   CountryStatus = "active"
	CountryStatusInactive CountryStatus = "inactive"
)

type Country struct {
	ID        int64         `db:"id"`
	Name      string        `db:"name"`
	Code      string        `db:"code"`
	Currency  string        `db:"currency"`
	PhoneCode string        `db:"phone_code"`
	TimeZone  string        `db:"time_zone"`
	Status    CountryStatus `db:"status"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
}
