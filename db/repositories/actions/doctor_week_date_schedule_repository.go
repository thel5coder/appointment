package actions

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/db/repositories/contracts"
)

type DoctorWeekDateScheduleRepository struct {
	DB *sql.DB
}

func NewDoctorWeekDateScheduleRepository(DB *sql.DB) contracts.IDoctorWeekDateScheduleRepository {
	return &DoctorWeekDateScheduleRepository{DB: DB}
}

func (DoctorWeekDateScheduleRepository) Add(model models.DoctorWeekDateSchedule, tx *sql.Tx) (err error) {
	statement := `insert into doctor_week_date_schedules (week_day_schedule_id,week_date,is_present,created_at,updated_at) values($1,$2,$3,$4,$5)`
	if _, err = tx.Exec(statement, model.WeekDayScheduleID, model.WeekDate, model.IsPresent, model.CreatedAt, model.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (DoctorWeekDateScheduleRepository) DeleteBy(column, value, operator string, model models.DoctorWeekDateSchedule, tx *sql.Tx) (err error) {
	statement := `update doctor_week_date_schedules set updated_at=$1, deleted_at=$2 where ` + column + `` + operator + `$3`
	if _, err = tx.Exec(statement, model.UpdatedAt, model.DeletedAt.Time, value); err != nil {
		return err
	}

	return nil
}
