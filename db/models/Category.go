package models

import (
	"database/sql"
	"time"
)

type Category struct {
	ID                 string         `db:"id"`
	Slug               string         `db:"slug"`
	Name               string         `db:"name"`
	ParentID           sql.NullString `db:"parent_id"`
	IsActive           bool           `db:"is_active"`
	FileIconID         sql.NullString `db:"icon"`
	FileIconPath       sql.NullString `db:"file_icon_path"`
	FileBackgroundID   sql.NullString `db:"file_background_id"`
	FileBackgroundPath sql.NullString `db:"file_background_path"`
	Treatments         sql.NullString `db:"treatments"`
	CreatedAt          time.Time      `db:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at"`
	DeletedAt          sql.NullTime   `db:"deleted_at"`
}
