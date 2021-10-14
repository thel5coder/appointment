package models

import (
	"database/sql"
	"time"
)

type DoctorWeekDayScheduleWorkTimes struct {
	ID                string       `db:"id"`
	WeekDayScheduleID string       `db:"week_day_schedule_id"`
	StartAt           time.Time    `db:"start_at"`
	EndAt             time.Time    `db:"end_at"`
	CreatedAt         time.Time    `db:"created_at"`
	UpdatedAt         time.Time    `db:"updated_at"`
	DeletedAt         sql.NullTime `db:"deleted_at"`
}
