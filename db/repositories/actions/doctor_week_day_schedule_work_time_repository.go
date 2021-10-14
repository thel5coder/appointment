package actions

import (
	"database/sql"
	"fmt"
	"profira-backend/db/models"
	"profira-backend/db/repositories/contracts"
)

type DoctorWorkDayScheduleWorkTimeRepository struct {
	DB *sql.DB
}

func NewDoctorWorkDayScheduleWorkTimeRepository(DB *sql.DB) contracts.IDoctorWeekDayScheduleWorkTimeRepository {
	return &DoctorWorkDayScheduleWorkTimeRepository{DB: DB}
}

func (DoctorWorkDayScheduleWorkTimeRepository) Add(model models.DoctorWeekDayScheduleWorkTimes, tx *sql.Tx) (err error) {
	statement := `insert into doctor_weekday_schedule_work_times (week_day_schedule_id,start_at,end_at,created_at,updated_at) values($1,$2,$3,$4,$5)`
	if _, err = tx.Exec(statement, model.WeekDayScheduleID, model.StartAt, model.EndAt, model.CreatedAt, model.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (DoctorWorkDayScheduleWorkTimeRepository) Edit(model models.DoctorWeekDayScheduleWorkTimes, tx *sql.Tx) (err error) {
	statement := `update doctor_weekday_schedule_work_times set start_at=$1, end_at=$2, updated_at=$3 where id=$4`
	if _, err = tx.Exec(statement, model.StartAt, model.EndAt, model.UpdatedAt, model.ID); err != nil {
		return err
	}

	return nil
}

func (DoctorWorkDayScheduleWorkTimeRepository) DeleteBy(column, value, operator string, model models.DoctorWeekDayScheduleWorkTimes, tx *sql.Tx) (err error) {
	fmt.Println(model.UpdatedAt)
	statement := `update doctor_weekday_schedule_work_times set updated_at=$1, deleted_at=$2 where ` + column + `` + operator + `$3`
	if _, err = tx.Exec(statement, model.UpdatedAt, model.DeletedAt.Time, value); err != nil {
		return err
	}

	return nil
}
