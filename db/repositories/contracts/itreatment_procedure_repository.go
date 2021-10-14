package contracts

import (
	"database/sql"
	"profira-backend/db/models"
)

type ITreatmentProcedureRepository interface {
	Edit(model models.TreatmentProcedure,tx *sql.Tx) (err error)

	Add(model models.TreatmentProcedure,tx *sql.Tx) (err error)

	DeleteBy(column,value,operator string,model models.TreatmentProcedure,tx *sql.Tx) (err error)
}
