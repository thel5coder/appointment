package contracts

import (
	"database/sql"
	"profira-backend/db/models"
)

type IBedRepository interface {
	Browse(clinicID, search, order, sort string, limit, offset int) (data []models.Bed, count int, err error)

	BrowseAll(clinicID, search string) (data []models.Bed, err error)

	ReadBy(column, value, operator string) (data models.Bed, err error)

	Edit(input models.Bed, tx *sql.Tx) (err error)

	Add(input models.Bed,tx *sql.Tx) (res string, err error)

	DeleteBy(column, value, operator string, input models.Bed,tx *sql.Tx) (err error)

	CountBy(ID, column, value, operator string) (res int, err error)
}
