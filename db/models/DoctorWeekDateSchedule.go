package models

import (
	"database/sql"
	"time"
)

type DoctorWeekDateSchedule struct {
	ID                string       `db:"id"`
	WeekDayScheduleID string       `db:"week_day_schedule_id"`
	WeekDate          time.Time    `db:"week_date"`
	IsPresent         bool         `db:"is_present"`
	CreatedAt         time.Time    `db:"created_at"`
	UpdatedAt         time.Time    `db:"updated_at"`
	DeletedAt         sql.NullTime `db:"deleted_at"`
}
