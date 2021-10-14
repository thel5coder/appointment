package contracts

import "profira-backend/db/models"

type ICategoryRepository interface {
	Browse(search, order, sort string, limit, offset int) (data []models.Category,count int,err error)

	BrowseTree(search string) (data []models.Category, err error)

	BrowseByParentID(parentID string) (data []models.Category, err error)

	ReadBy(column, value, operator string) (data models.Category, err error)

	Edit(model models.Category) (res string, err error)

	Add(model models.Category) (res string, err error)

	Delete(model models.Category) (res string, err error)
}
