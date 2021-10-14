package contracts

import (
	"database/sql"
	"profira-backend/db/models"
)

type IDoctorWeekDayScheduleWorkTimeRepository interface {
	Add(model models.DoctorWeekDayScheduleWorkTimes,tx *sql.Tx) (err error)

	Edit(model models.DoctorWeekDayScheduleWorkTimes,tx *sql.Tx) (err error)

	DeleteBy(column,value,operator string,model models.DoctorWeekDayScheduleWorkTimes,tx *sql.Tx) (err error)
}
