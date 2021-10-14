package contracts

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/usecase/viewmodel"
)

type IStaffRepository interface {
	BrowseByRole(roleID, search, order, sort string, limit, offset int) (data []models.Staff, count int, err error)

	BrowseAllByRole(roleID, search string) (data []models.Staff, err error)

	BrowseStaffDoctorByClinic(clinicId,roleId string) (data []models.Staff,err error)

	ReadBy(column, value, operator string) (data models.Staff, err error)

	Edit(input viewmodel.StaffVm, tx *sql.Tx) (err error)

	Add(input viewmodel.StaffVm, tx *sql.Tx) (res string, err error)

	DeleteBy(column, value, operator, updatedAt, deletedAt string, tx *sql.Tx) (err error)

	CountBy(ID, column, value, operator string) (res int, err error)
}
