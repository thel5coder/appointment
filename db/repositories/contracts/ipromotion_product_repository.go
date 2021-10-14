package contracts

import (
	"database/sql"
	"profira-backend/db/models"
)

type IPromotionProductRepository interface {
	Add(model models.PromotionProduct, tx *sql.Tx) (err error)

	DeleteBy(column, value, operator string, model models.PromotionProduct, tx *sql.Tx) (err error)

	CountBy(column,value,operator string) (res int,err error)
}
