package actions

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/db/repositories/contracts"
	"strings"
)

type PromotionRepository struct {
	DB *sql.DB
}

func NewPromotionRepository(DB *sql.DB) contracts.IPromotionRepository {
	return &PromotionRepository{DB: DB}
}

const (
	promotionSelectStatement = `select p.id,p.slug,p.name,p.customer_promotion_condition,p.promotion_type,p.description,p.start_date,p.end_date,p.foto_id,f.path,
                                p.nominal_type, p.nominal_percentage, p.nominal_amount, p.birth_date_condition, p.sex_condition, p.register_date_condition_start, 
                                p.register_date_condition_end, p.created_at,p.updated_at,array_to_string(array_agg(mt.id ||':'|| mt.name),',')`
	promotionJoinStatement = `left join files f on f.id = p.foto_id and f.deleted_at is null
                                 left join promotion_products pp on pp.promotion_id = p.id and pp.deleted_at is null
                                 left join master_treatments mt on mt.id = pp.product_id and mt.deleted_at is null`
	promotionWhereStatement   = `where (lower(p.name) like $1) and p.deleted_at is null`
	promotionGroupByStatement = `group by p.id,f.id`
)

func (repository PromotionRepository) scanRow(row *sql.Row) (res models.Promotion, err error) {
	err = row.Scan(&res.ID, &res.Slug, &res.Name, &res.CustomerPromotionCondition, &res.PromotionType, &res.Description, &res.StartDate, &res.EndDate, &res.FotoID, &res.FilePath,
		&res.NominalType, &res.NominalPercentage, &res.NominalAmount, &res.BirthDateCondition, &res.SexCondition, &res.RegisterDateConditionStart,
		&res.RegisterDateConditionEnd, &res.CreatedAt, &res.UpdatedAt, &res.Treatments)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository PromotionRepository) scanRows(rows *sql.Rows) (res models.Promotion, err error) {
	err = rows.Scan(&res.ID, &res.Slug, &res.Name, &res.CustomerPromotionCondition, &res.PromotionType, &res.Description, &res.StartDate, &res.EndDate, &res.FotoID, &res.FilePath,
		&res.NominalType, &res.NominalPercentage, &res.NominalAmount, &res.BirthDateCondition, &res.SexCondition, &res.RegisterDateConditionStart,
		&res.RegisterDateConditionEnd, &res.CreatedAt, &res.UpdatedAt, &res.Treatments)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository PromotionRepository) Browse(search, orderBy, sort string, limit, offset int) (res []models.Promotion, count int, err error) {
	statement := promotionSelectStatement + ` from promotions p ` + promotionJoinStatement + ` ` + promotionWhereStatement + ` ` + promotionGroupByStatement +
		` order by ` + orderBy + ` ` + sort + ` limit $2 offset $3`
	rows, err := repository.DB.Query(statement, "%"+strings.ToLower(search)+"%", limit, offset)
	if err != nil {
		return res, count, err
	}

	for rows.Next() {
		temp, err := repository.scanRows(rows)
		if err != nil {
			return res, count, err
		}
		res = append(res, temp)
	}

	statement = `select count(distinct p.id) from promotions p ` + promotionJoinStatement + ` ` + promotionWhereStatement + ` ` + promotionGroupByStatement
	err = repository.DB.QueryRow(statement, "%"+strings.ToLower(search)+"%").Scan(&count)
	if err != nil {
		return res, count, err
	}

	return res, count, nil
}

func (repository PromotionRepository) BrowseActivePromotion(startAt, endAt int64) (res []models.Promotion, err error) {
	statement := promotionSelectStatement + ` from promotions p ` + promotionJoinStatement + ` where p.deleted_at is null ` + promotionGroupByStatement
	rows, err := repository.DB.Query(statement)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		temp, err := repository.scanRows(rows)
		if err != nil {
			return res, err
		}
		res = append(res, temp)
	}

	return res, nil
}

func (repository PromotionRepository) ReadBy(column, value, operator string) (res models.Promotion, err error) {
	statement := promotionSelectStatement + ` from promotions p ` + promotionJoinStatement + ` where ` + column + `` + operator + `$1 and p.deleted_at is null ` + promotionGroupByStatement
	row := repository.DB.QueryRow(statement, value)
	res, err = repository.scanRow(row)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (PromotionRepository) Edit(model models.Promotion, tx *sql.Tx) (err error) {
	statement := `update promotions set name=$1, slug=$2, customer_promotion_condition=$3, promotion_type=$4, description=$5, start_date=$6, end_date=$7,
                  foto_id=$8, nominal_type=$9, nominal_percentage=$10, nominal_amount=$11, birth_date_condition=$12, sex_condition=$13, register_date_condition_start=$14, register_date_condition_end=$15,
                  updated_at=$16,start_at_unix=$17, end_at_unix=$18 where id=$19`
	_, err = tx.Exec(statement, model.Name, model.Slug, model.CustomerPromotionCondition.String, model.PromotionType, model.Description, model.StartDate, model.EndDate,
		model.FotoID.String, model.NominalType, model.NominalPercentage.Int32, model.NominalAmount.Int32, model.BirthDateCondition.Time, model.SexCondition.String,
		model.RegisterDateConditionStart.Time, model.RegisterDateConditionEnd.Time, model.UpdatedAt, model.StartAtUnix, model.EndAtUnix, model.ID)
	if err != nil {
		return err
	}

	return nil
}

func (PromotionRepository) Add(model models.Promotion, tx *sql.Tx) (res string, err error) {
	statement := `insert into promotions (name,slug,customer_promotion_condition,promotion_type,description,start_date,end_date,foto_id,nominal_type, nominal_percentage,
                 nominal_amount, birth_date_condition, sex_condition, register_date_condition_start, register_date_condition_end,created_at,updated_at,start_at_unix,end_at_unix) 
                 values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) returning id`
	err = tx.QueryRow(statement, model.Name, model.Slug, model.CustomerPromotionCondition.String, model.PromotionType, model.Description, model.StartDate, model.EndDate, model.FotoID.String,
		model.NominalType, model.NominalPercentage.Int32, model.NominalAmount.Int32, model.BirthDateCondition.Time, model.SexCondition.String, model.RegisterDateConditionStart.Time, model.RegisterDateConditionEnd.Time,
		model.CreatedAt, model.UpdatedAt, model.StartAtUnix, model.EndAtUnix).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (PromotionRepository) Delete(model models.Promotion, tx *sql.Tx) (err error) {
	statement := `update promotions set updated_at=$1, deleted_at=$2 where id=$3`
	_, err = tx.Exec(statement, model.UpdatedAt, model.DeletedAt.Time,model.ID)
	if err != nil {
		return err
	}

	return nil
}
