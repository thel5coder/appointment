package models

import "database/sql"

type Docter struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	UserID      string         `db:"user_id"`
	User        User           `db:"user"`
	Clinics     string         `db:"clinics"`
	Treatment   sql.NullString `db:"treatment"`
	CreatedAt   string         `db:"created_at"`
	UpdatedAt   string         `db:"updated_at"`
}
