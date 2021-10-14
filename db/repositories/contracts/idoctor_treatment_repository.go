package contracts

import (
	"database/sql"
	"profira-backend/db/models"
)

type IDoctorTreatmentRepository interface {
	Add(model models.DoctorTreatment,tx *sql.Tx) (err error)

	DeleteBy(column,value,operator string,model models.DoctorTreatment,tx *sql.Tx) (err error)

	CountBy(column,value,operator string) (res int,err error)
}
