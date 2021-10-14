package actions

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/db/repositories/contracts"
)

type PromotionProductRepository struct {
	DB *sql.DB
}

func NewPromotionProductRepository(db *sql.DB) contracts.IPromotionProductRepository {
	return PromotionProductRepository{DB: db}
}

func (PromotionProductRepository) Add(model models.PromotionProduct, tx *sql.Tx) (err error) {
	statement := `insert into promotion_products (product_id,promotion_id,created_at,updated_at) values($1,$2,$3,$4)`
	_, err = tx.Exec(statement, model.ProductID, model.PromotionID, model.CreatedAt, model.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (PromotionProductRepository) DeleteBy(column, value, operator string, model models.PromotionProduct, tx *sql.Tx) (err error) {
	statement := `update promotion_products set updated_at=$1, deleted_at=$2 where ` + column + `` + operator + `$3`
	_, err = tx.Exec(statement, model.UpdatedAt, model.DeletedAt.Time, value)
	if err != nil {
		return err
	}

	return nil
}

func (repository PromotionProductRepository) CountBy(column, value, operator string) (res int, err error) {
	statement := `select count(id) from promotion_products where ` + column + `` + operator + `$1 and deleted_at is null`
	err = repository.DB.QueryRow(statement, value).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}
