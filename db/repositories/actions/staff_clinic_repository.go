package actions

import (
	"database/sql"
	"profira-backend/db/repositories/contracts"
	"profira-backend/helpers/datetime"
	"time"
)

type StaffClinicRepository struct {
	DB *sql.DB
}

func NewStaffClinicRepository(DB *sql.DB) contracts.IStaffClinicRepository {
	return &StaffClinicRepository{DB: DB}
}

func (StaffClinicRepository) Add(staffID, clinicID, createdAt string, tx *sql.Tx) (err error) {
	statement := `insert into "staff_clinics" ("staff_id","clinic_id","created_at") values($1,$2,$3)`
	_, err = tx.Exec(statement, staffID, clinicID, datetime.StrParseToTime(createdAt, time.RFC3339))
	if err != nil {
		return err
	}

	return nil
}

func (StaffClinicRepository) DeleteBy(column, value, operator string, tx *sql.Tx) (err error) {
	statement := `delete from "staff_clinics" where ` + column + `` + operator + `$1`
	_, err = tx.Exec(statement, value)
	if err != nil {
		return err
	}

	return nil
}

func (repository StaffClinicRepository) CountBy(column, value, operator string) (res int, err error) {
	statement := `select count("id") from "staff_clinics" where ` + column + `` + operator + `$1`
	err = repository.DB.QueryRow(statement, value).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}
