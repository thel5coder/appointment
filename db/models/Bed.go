package models

import (
	"database/sql"
	"time"
)

type Bed struct {
	ID         string         `db:"id"`
	BedCode    string         `db:"bed_code"`
	ClinicID   string         `db:"clinic_id"`
	IsUseAble  bool           `db:"is_use_able"`
	CreatedAt  time.Time      `db:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at"`
	DeletedAt  sql.NullTime   `db:"deleted_at"`
	Clinic     Clinic         `db:"clinic"`
	Treatments sql.NullString `db:"treatments"`
}
