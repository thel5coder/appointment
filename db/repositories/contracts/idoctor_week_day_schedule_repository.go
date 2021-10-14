package contracts

import (
	"database/sql"
	"profira-backend/db/models"
)

type IDoctorWeekDayScheduleRepository interface {
	BrowseBy(filters map[string]interface{}) (res []models.DoctorWeekDaySchedule,err error)

	Add(model models.DoctorWeekDaySchedule,tx *sql.Tx) (res string,err error)

	DeleteBy(column,value,operator string,model models.DoctorWeekDaySchedule,tx *sql.Tx) (err error)

	CountBy(column,value,operator string) (res int,err error)
}
