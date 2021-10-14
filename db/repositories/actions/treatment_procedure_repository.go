package actions

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/db/repositories/contracts"
)

type TreatmentProcedureRepository struct {
	DB *sql.DB
}

func NewTreatmentProcedureRepository(DB *sql.DB) contracts.ITreatmentProcedureRepository {
	return &TreatmentProcedureRepository{DB: DB}
}

func (TreatmentProcedureRepository) Edit(model models.TreatmentProcedure, tx *sql.Tx) (err error) {
	statement := `update "treatment_procedures" set procedure_id=$1, updated_at=$2 where id=$3`
	_, err = tx.Exec(statement, model.ProcedureID, model.UpdatedAt, model.ID)
	if err != nil {
		return err
	}

	return nil
}

func (TreatmentProcedureRepository) Add(model models.TreatmentProcedure, tx *sql.Tx) (err error) {
	statement := `insert into "treatment_procedures" (procedure_id,treatment_id,created_at,updated_at) values($1,$2,$3,$4)`
	_, err = tx.Exec(statement, model.ProcedureID, model.TreatmentID, model.CreatedAt, model.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (TreatmentProcedureRepository) DeleteBy(column, value, operator string, model models.TreatmentProcedure, tx *sql.Tx) (err error) {
	statement := `update "treatment_procedures" set updated_at=$1, deleted_at=$2 where ` + column + `` + operator + `$3`
	_, err = tx.Exec(statement, model.UpdatedAt, model.DeletedAt.Time, value)
	if err != nil {
		return err
	}

	return nil
}
