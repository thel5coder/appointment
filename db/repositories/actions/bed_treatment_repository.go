package actions

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/db/repositories/contracts"
)

type BedTreatmentRepository struct{
	DB *sql.DB
}

func NewBedTreatmentRepository(DB *sql.DB) contracts.IBedTreatmentRepository{
	return &BedTreatmentRepository{DB: DB}
}

func (BedTreatmentRepository) Add(model models.BedTreatment, tx *sql.Tx) (err error) {
	statement := `insert into bed_treatments (bed_id,treatment_id,created_at,updated_at) values($1,$2,$3,$4)`
	_,err = tx.Exec(statement,model.BedID,model.TreatmentID,model.CreatedAt,model.UpdatedAt)
	if err != nil{
		return err
	}

	return nil
}

func (BedTreatmentRepository) DeleteBy(column, value, operator string, model models.BedTreatment, tx *sql.Tx) (err error) {
	statement := `update bed_treatments set updated_at=$1, deleted_at=$2 where `+column+``+operator+`$3`
	_,err = tx.Exec(statement,model.UpdatedAt,model.DeletedAt.Time,value)
	if err != nil {
		return err
	}

	return nil
}

func (repository BedTreatmentRepository) CountBy(column, value, operator string) (res int, err error) {
	statement := `select count(id) from bed_treatments where `+column+``+operator+`$1 and deleted_at is null`
	err = repository.DB.QueryRow(statement,value).Scan(&res)
	if err != nil {
		return res,err
	}

	return res,nil
}

