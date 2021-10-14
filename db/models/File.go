package models

import (
	"database/sql"
	"time"
)

type File struct {
	ID        string       `json:"id"`
	Path      string       `json:"path"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}
