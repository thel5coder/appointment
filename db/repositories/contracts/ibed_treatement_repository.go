package contracts

import (
	"database/sql"
	"profira-backend/db/models"
)

type IBedTreatmentRepository interface {
	Add(model models.BedTreatment,tx *sql.Tx) (err error)

	DeleteBy(column,value,operator string,model models.BedTreatment,tx *sql.Tx) (err error)

	CountBy(column,value,operator string) (res int,err error)
}
