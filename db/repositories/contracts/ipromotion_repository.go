package contracts

import (
	"database/sql"
	"profira-backend/db/models"
)

type IPromotionRepository interface {
	Browse(search,orderBy,sort string,limit,offset int) (res []models.Promotion,count int,err error)

	BrowseActivePromotion(startAt,endAt int64) (res []models.Promotion,err error)

	ReadBy(column,value,operator string) (res models.Promotion,err error)

	Edit(model models.Promotion,tx *sql.Tx) (err error)

	Add(model models.Promotion,tx *sql.Tx) (res string,err error)

	Delete(model models.Promotion,tx *sql.Tx) (err error)
}
