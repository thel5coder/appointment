package contracts

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/usecase/viewmodel"
)

type IUserRepository interface {
	Browse(search, order, sort string, limit, offset int) (data []models.User, count int, err error)

	BrowseAllBy(column, value string) (data []models.User, err error)

	ReadBy(column, value, operator string) (data models.User, err error)

	Edit(input viewmodel.UserVm, password string, tx *sql.Tx) (err error)

	EditPassword(ID, password, updatedAt string) (res string, err error)

	EditFcmDeviceToken(ID, fcmDeviceToken, updatedAt string) (res string, err error)

	EditActivatedUser(ID,updatedAt,activatedAt string,isActive bool) (res string,err error)

	Add(input viewmodel.UserVm, password string, tx *sql.Tx) (res string, err error)

	DeleteBy(column, value, operator, updatedAt, deletedAt string, tx *sql.Tx) (err error)

	CountBy(ID, column, value, operator string) (res int, err error)
}
