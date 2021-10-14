package actions

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/db/repositories/contracts"
)

type DoctorTreatmentRepository struct {
	DB *sql.DB
}

func NewDoctorTreatmentRepository(DB *sql.DB) contracts.IDoctorTreatmentRepository {
	return &DoctorTreatmentRepository{DB: DB}
}

func (DoctorTreatmentRepository) Add(model models.DoctorTreatment, tx *sql.Tx) (err error) {
	statement := `insert into doctor_treatments (doctor_id,treatment_id,created_at,updated_at) values($1,$2,$3,$4)`
	_, err = tx.Exec(statement, model.DoctorID, model.TreatmentID, model.CreatedAt, model.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (DoctorTreatmentRepository) DeleteBy(column, value, operator string, model models.DoctorTreatment, tx *sql.Tx) (err error) {
	statement := `update doctor_treatments set updated_at=$1, deleted_at=$2 where ` + column + `` + operator + `$3`
	_, err = tx.Exec(statement, model.UpdatedAt, model.DeletedAt.Time, value)
	if err != nil {
		return err
	}

	return nil
}

func (repository DoctorTreatmentRepository) CountBy(column, value, operator string) (res int, err error) {
	statement := `select count(id) from doctor_treatments where ` + column + `` + operator + `$1 and deleted_at is null`
	err = repository.DB.QueryRow(statement, value).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}
