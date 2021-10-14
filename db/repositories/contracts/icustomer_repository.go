package contracts

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/usecase/viewmodel"
)

type ICustomerRepository interface {
	Browse(search, order, sort string, limit, offset int) (data []models.Customer, count int, err error)

	BrowseAllBy(column, value string) (data []models.Customer, err error)

	ReadBy(column, value, operator string) (data models.Customer, err error)

	Edit(input viewmodel.CustomerVm, tx *sql.Tx) (err error)

	EditProfile(model models.Customer, tx *sql.Tx) (err error)

	Add(input viewmodel.CustomerVm, tx *sql.Tx) (err error)

	DeleteBy(column, value, operator, updatedAt, deletedAt string, tx *sql.Tx) (err error)

	CountBy(ID, column, value, operator string) (res int, err error)
}
