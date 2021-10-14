package models

import "database/sql"

type Staff struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	UserID      string         `db:"user_id"`
	User        User           `db:"user"`
	Clinics     string         `db:"clinics"`
	Treatments  sql.NullString `db:"treatments"`
	CreatedAt   string         `db:"created_at"`
	UpdatedAt   string         `db:"updated_at"`
}
