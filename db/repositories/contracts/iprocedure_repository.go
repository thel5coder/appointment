package contracts

import (
	"profira-backend/db/models"
)

type IProcedureRepository interface {
	Browse(search, order, sort string, limit, offset int) (data []models.Procedure, count int, err error)

	BrowseAll(search string) (data []models.Procedure,err error)

	ReadBy(column, value,operator string) (data models.Procedure, err error)

	Edit(model models.Procedure) (res string, err error)

	Add(model models.Procedure) (res string, err error)

	DeleteBy(column,value,operator string, model models.Procedure) (res string, err error)

	CountBy(ID, column, value,operator string) (res int, err error)
}
