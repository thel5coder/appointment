package contracts

import (
	"profira-backend/db/models"
	"profira-backend/usecase/viewmodel"
)

type IRoleRepository interface {
	Browse(search, order, sort string, limit, offset int) (data []models.Role, count int, err error)

	ReadBy(column, value string) (data models.Role, err error)

	Edit(input viewmodel.RoleVm) (res string, err error)

	Add(input viewmodel.RoleVm) (res string, err error)

	Delete(ID, updatedAt, deletedAt string) (res string, err error)

	CountBy(ID, column, value string) (res int, err error)
}
