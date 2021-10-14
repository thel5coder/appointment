package models

import (
	"database/sql"
	"time"
)

type TreatmentProcedure struct {
	ID          string       `db:"id"`
	ProcedureID string       `db:"procedure_id"`
	TreatmentID string       `db:"treatment_id"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}
