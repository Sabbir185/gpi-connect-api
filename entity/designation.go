package entity

import "time"

type Designation struct {
	ID        int64     `db:"id"`
	Title     string    `db:"title"`
	Hint      *string   `db:"hint"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
