package models

import "database/sql"

type Role struct {
	ID        string         `db:"id"`
	Name      string         `db:"name"`
	Slug      string         `db:"slug"`
	CreatedAt string         `db:"created_at"`
	UpdatedAt string         `db:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at"`
}
