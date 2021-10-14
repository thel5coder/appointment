package actions

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/db/repositories/contracts"
)

type DoctorWeekDayScheduleRepository struct {
	DB *sql.DB
}

func NewDoctorWeekDayScheduleRepository(DB *sql.DB) contracts.IDoctorWeekDayScheduleRepository {
	return &DoctorWeekDayScheduleRepository{DB: DB}
}

const (
	doctorWeekdayScheduleSelectStatement = `select array_to_string(array_agg(dw.id ||'|'|| dw.day ||'|'|| dw.created_at ||'|'|| dw.updated_at),','),ms.id,ms.name,mc.id,mc.name,
                                            array_to_string(array_agg(dwt.id ||'|'|| dw.id ||'|'|| dwt.start_at || '|' || dwt.end_at),',')`
	doctorWeekdayScheduleJoinStatement = `inner join master_staffs ms on ms.id = dw.doctor_id and ms.deleted_at is null
                                          inner join master_clinics mc on mc.id = dw.clinic_id and mc.deleted_at is null
                                          left join doctor_weekday_schedule_work_times dwt on dwt.week_day_schedule_id = dw.id and dwt.deleted_at is null`
	doctorWeekdayGroupByStatement = `group by ms.id,mc.id`
)

var (
	whereStatement = `where dw.deleted_at is null`
)

func (repository DoctorWeekDayScheduleRepository) scanRows(rows *sql.Rows) (res models.DoctorWeekDaySchedule, err error) {
	if err = rows.Scan(&res.ScheduleDays, &res.DoctorID, &res.Staff.Name, &res.ClinicID, &res.Clinic.Name, &res.ScheduleTimes); err != nil {
		return res, err
	}

	return res, nil
}

func (repository DoctorWeekDayScheduleRepository) scanRow(row *sql.Rows) (res models.DoctorWeekDaySchedule, err error) {
	if err = row.Scan(&res.ScheduleDays, &res.DoctorID, &res.Staff.Name, &res.ClinicID, &res.Clinic.Name, &res.ScheduleTimes); err != nil {
		return res, err
	}

	return res, nil
}

func (repository DoctorWeekDayScheduleRepository) BrowseBy(filters map[string]interface{}) (res []models.DoctorWeekDaySchedule, err error) {
	whereConditionStatement := ``
	if val, ok := filters["clinic_id"]; ok {
		whereConditionStatement += ` and dw.clinic_id='` + val.(string) + `'`
	}
	if val, ok := filters["doctor_id"]; ok {
		whereConditionStatement += ` and dw.doctor_id='` + val.(string) + `'`
	}

	statement := doctorWeekdayScheduleSelectStatement + ` from doctor_week_day_schedules dw ` + doctorWeekdayScheduleJoinStatement + ` ` + whereStatement + ` ` + whereConditionStatement +
		` ` + doctorWeekdayGroupByStatement
	rows, err := repository.DB.Query(statement)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		temp, err := repository.scanRows(rows)
		if err != nil {
			return res, err
		}
		res = append(res, temp)
	}

	return res, nil
}

func (DoctorWeekDayScheduleRepository) Add(model models.DoctorWeekDaySchedule, tx *sql.Tx) (res string, err error) {
	statement := `insert into doctor_week_day_schedules (doctor_id,clinic_id,day,created_at,updated_at) values($1,$2,$3,$4,$5) returning id`
	if err = tx.QueryRow(statement, model.DoctorID, model.ClinicID, model.Day, model.CreatedAt, model.UpdatedAt).Scan(&res); err != nil {
		return res, err
	}

	return res, nil
}

func (DoctorWeekDayScheduleRepository) DeleteBy(column, value, operator string, model models.DoctorWeekDaySchedule, tx *sql.Tx) (err error) {
	statement := `update doctor_week_day_schedules set updated_at=$1, deleted_at=$2 where ` + column + `` + operator + `$3`
	if _, err = tx.Exec(statement, model.UpdatedAt, model.DeletedAt.Time, value); err != nil {
		return err
	}

	return nil
}

func (repository DoctorWeekDayScheduleRepository) CountBy(column, value, operator string) (res int, err error) {
	statement := `select count(distinct dw.id) from doctor_week_day_schedules where ` + column + `` + operator + `$1 and deleted_at is null`
	if err = repository.DB.QueryRow(statement, value).Scan(&res); err != nil {
		return res, err
	}

	return res, nil
}
