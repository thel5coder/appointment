package contracts

import "database/sql"

type IStaffClinicRepository interface {
	Add(staffID, clinicID, createdAt string, tx *sql.Tx) (err error)

	DeleteBy(column, value, operator string, tx *sql.Tx) (err error)

	CountBy(column, value, operator string) (res int, err error)
}
