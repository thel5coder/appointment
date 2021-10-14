package models

import (
	"database/sql"
	"time"
)

type DoctorTreatment struct {
	ID          string       `db:"id"`
	DoctorID    string       `db:"doctor_id"`
	TreatmentID string       `db:"treatment_id"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}
