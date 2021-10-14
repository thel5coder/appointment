package models

import (
	"database/sql"
	"time"
)

type DoctorWeekDaySchedule struct {
	ID            string         `db:"id"`
	DoctorID      string         `db:"doctor_id"`
	ClinicID      string         `db:"clinic_id"`
	Day           string         `db:"day"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
	DeletedAt     sql.NullTime   `db:"deleted_at"`
	Staff         Staff          `db:"staff"`
	Clinic        Clinic         `db:"clinic"`
	ScheduleTimes sql.NullString `db:"schedule_times"`
	ScheduleDays  sql.NullString `db:"schedule_days"`
}
