package contracts

import (
	"database/sql"
	"profira-backend/db/models"
)

type IDoctorWeekDateScheduleRepository interface {
	Add(model models.DoctorWeekDateSchedule,tx *sql.Tx) (err error)

	DeleteBy(column,value,operator string,model models.DoctorWeekDateSchedule,tx *sql.Tx) (err error)
}
