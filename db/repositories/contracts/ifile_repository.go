package contracts

import (
	"profira-backend/db/models"
)

type IFileRepository interface {
	ReadBy(column, value, operator string) (data models.File, err error)

	Add(model models.File) (res string, err error)

	Delete(model models.File) (res string, err error)

	CountBy(column, value string) (res int, err error)
}
