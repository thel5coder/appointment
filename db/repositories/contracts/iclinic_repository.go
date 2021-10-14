package contracts

import (
	"profira-backend/db/models"
	"profira-backend/usecase/viewmodel"
)

type IClinicRepository interface {
	Browse(search, order, sort string, limit, offset int) (data []models.Clinic, count int, err error)

	BrowseAll(search string) (data []models.Clinic, err error)

	ReadBy(column, value, operator string) (data models.Clinic, err error)

	Edit(input viewmodel.ClinicVm) (res string, err error)

	Add(input viewmodel.ClinicVm) (res string, err error)

	Delete(ID, updatedAt, deletedAt string) (res string, err error)

	CountBy(ID, column, value, operator string) (res int, err error)
}
