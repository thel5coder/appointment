package models

import (
	"database/sql"
	"time"
)

type BedTreatment struct {
	ID          string       `json:"id"`
	BedID       string       `json:"bed_id"`
	TreatmentID string       `json:"treatment_id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
}
