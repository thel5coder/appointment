package models

import (
	"database/sql"
	"time"
)

type PromotionProduct struct {
	ID          string       `db:"id"`
	ProductID   string       `db:"product_id"`
	PromotionID string       `db:"promotion_id"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}
