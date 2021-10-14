package models

import (
	"database/sql"
	"time"
)

type Procedure struct {
	ID        string       `db:"id"`
	Name      string       `db:"name"`
	Duration  int          `db:"duration"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
