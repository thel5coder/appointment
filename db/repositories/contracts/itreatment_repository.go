package contracts

import (
	"database/sql"
	"profira-backend/db/models"
)

type ITreatmentRepository interface {
	Browse(search, order, sort string, limit, offset int) (data []models.Treatment, count int, err error)

	BrowseByCategory(categoryID string) (data []models.Treatment,err error)

	BrowseAll(search string) (data []models.Treatment,err error)

	ReadBy(column, value, operator string) (data models.Treatment, err error)

	Edit(model models.Treatment, tx *sql.Tx) (err error)

	Add(model models.Treatment, tx *sql.Tx) (res string, err error)

	Delete(model models.Treatment, tx *sql.Tx) (err error)

	CountBy(ID,column,value,categoryId,operator string) (res int,err error)
}
