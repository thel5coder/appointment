package models

import (
	"database/sql"
	"time"
)

type Treatment struct {
	ID                  string         `db:"id"`
	Name                string         `db:"name"`
	Description         string         `db:"description"`
	Duration            int32          `db:"duration"`
	PhotoID             sql.NullString `db:"photo_id"`
	PhotoPath           sql.NullString `db:"photo_path"`
	IconID              sql.NullString `db:"icon_id"`
	IconPath            sql.NullString `db:"icon_path"`
	CreatedAt           time.Time      `db:"created_at"`
	UpdatedAt           time.Time      `db:"updated_at"`
	DeletedAt           sql.NullTime   `db:"deleted_at"`
	Price               sql.NullInt32  `db:"price"`
	TreatmentProcedures sql.NullString `db:"treatment_procedures"`
	CategoryID          string         `db:"category_id"`
	CategoryName        string         `db:"category_name"`
}
