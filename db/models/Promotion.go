package models

import (
	"database/sql"
	"time"
)

type Promotion struct {
	ID                         string         `db:"id"`
	Slug                       string         `db:"slug"`
	Name                       string         `db:"name"`
	CustomerPromotionCondition sql.NullString `db:"customer_promotion_condition"`
	PromotionType              string         `db:"promotion_type"`
	Description                string         `db:"description"`
	StartDate                  time.Time      `db:"start_date"`
	StartAtUnix                int64          `db:"start_at_unix"`
	EndAtUnix                  int64          `db:"end_at_unix"`
	EndDate                    time.Time      `db:"end_date"`
	FotoID                     sql.NullString `db:"foto_id"`
	FilePath                   sql.NullString `db:"file_path"`
	NominalType                string         `db:"nominal_type"`
	NominalPercentage          sql.NullInt32  `db:"nominal_percentage"`
	NominalAmount              sql.NullInt32  `db:"nominal_amount"`
	BirthDateCondition         sql.NullTime   `db:"birth_date_condition"`
	SexCondition               sql.NullString `db:"sex_condition"`
	RegisterDateConditionStart sql.NullTime   `db:"register_date_condition_start"`
	RegisterDateConditionEnd   sql.NullTime   `db:"register_date_condition_end"`
	Treatments                 sql.NullString `db:"treatments"`
	CreatedAt                  time.Time      `db:"created_at"`
	UpdatedAt                  time.Time      `db:"updated_at"`
	DeletedAt                  sql.NullTime   `db:"deleted_at"`
}
