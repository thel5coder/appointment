package actions

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/db/repositories/contracts"
	"strings"
)

type TreatmentRepository struct {
	DB *sql.DB
}

func NewTreatmentRepository(DB *sql.DB) contracts.ITreatmentRepository {
	return &TreatmentRepository{DB: DB}
}

const (
	treatmentSelectStatement = `select t.id,t.name,t.description,t.duration,t.price,t.photo_id,fp.path,t.icon_id,fi.path,t.created_at,t.updated_at,
                                 array_to_string(array_agg(tp.id || ':' || p.id || ':' || p.name || ':' || p.duration),','), t.category_id, c.name`
	treatmentJoinStatement = `left join "files" fp on fp.id = t.photo_id
                              left join "files" fi on fi.id = t.icon_id
                              left join "treatment_procedures" tp on tp.treatment_id=t.id and tp.deleted_at is null
                              left join "procedures" p on p.id = tp.procedure_id and p.deleted_at is null
                              left join categories c on c.id = t.category_id and c.deleted_at is null`
	treatmentWhereStatement   = `where (lower(t.name) like $1 or cast(t.duration as varchar) like $1 or lower(t.description) like $1) and t.deleted_at is null`
	treatmentGroupByStatement = `group by t.id,fp.id,fi.id,c.id`
)

func (repository TreatmentRepository) scanRows(rows *sql.Rows) (res models.Treatment, err error) {
	err = rows.Scan(&res.ID, &res.Name, &res.Description, &res.Duration, &res.Price, &res.PhotoID, &res.PhotoPath, &res.IconID, &res.IconPath, &res.CreatedAt, &res.UpdatedAt, &res.TreatmentProcedures, &res.CategoryID, &res.CategoryName)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository TreatmentRepository) scanRow(row *sql.Row) (res models.Treatment, err error) {
	err = row.Scan(&res.ID, &res.Name, &res.Description, &res.Duration, &res.Price, &res.PhotoID, &res.PhotoPath, &res.IconID, &res.IconPath, &res.CreatedAt, &res.UpdatedAt, &res.TreatmentProcedures, &res.CategoryID, &res.CategoryName)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository TreatmentRepository) Browse(search, order, sort string, limit, offset int) (data []models.Treatment, count int, err error) {
	statement := treatmentSelectStatement + ` from "master_treatments" t ` + treatmentJoinStatement + ` ` + treatmentWhereStatement + ` ` + treatmentGroupByStatement + ` order by ` + order + ` ` + sort + ` limit $2 offset $3`
	rows, err := repository.DB.Query(statement, "%"+strings.ToLower(search)+"%", limit, offset)
	if err != nil {
		return data, count, err
	}

	for rows.Next() {
		temp, err := repository.scanRows(rows)
		if err != nil {
			return data, count, err
		}
		data = append(data, temp)
	}

	statement = `select count(distinct t.id) from "master_treatments" t ` + treatmentJoinStatement + ` ` + treatmentWhereStatement
	err = repository.DB.QueryRow(statement, "%"+strings.ToLower(search)+"%").Scan(&count)
	if err != nil {
		return data, count, err
	}

	return data, count, nil
}

func (repository TreatmentRepository) BrowseByCategory(categoryID string) (data []models.Treatment, err error) {
	statement := treatmentSelectStatement + ` from master_treatments t ` + treatmentJoinStatement + ` where t.deleted_at is null and t.category_id=$1 ` + treatmentGroupByStatement+
		` ORDER BY t.name ASC`
	rows, err := repository.DB.Query(statement, categoryID)
	if err != nil {
		return data, err
	}

	for rows.Next() {
		temp, err := repository.scanRows(rows)
		if err != nil {
			return data, err
		}
		data = append(data, temp)
	}

	return data, nil
}

func (repository TreatmentRepository) BrowseAll(search string) (data []models.Treatment, err error) {
	statement := treatmentSelectStatement + ` from "master_treatments" t ` + treatmentJoinStatement + ` ` + treatmentWhereStatement + ` ` + treatmentGroupByStatement
	rows, err := repository.DB.Query(statement, "%"+strings.ToLower(search)+"%")
	if err != nil {
		return data, err
	}

	for rows.Next() {
		temp, err := repository.scanRows(rows)
		if err != nil {
			return data, err
		}
		data = append(data, temp)
	}

	return data, err
}

func (repository TreatmentRepository) ReadBy(column, value, operator string) (data models.Treatment, err error) {
	statement := treatmentSelectStatement + ` from "master_treatments" t ` + treatmentJoinStatement + ` where ` + column + `` + operator + `$1 and t.deleted_at is null ` + treatmentGroupByStatement
	row := repository.DB.QueryRow(statement, value)
	data, err = repository.scanRow(row)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (TreatmentRepository) Edit(model models.Treatment, tx *sql.Tx) (err error) {
	statement := `update "master_treatments" set name=$1, description=$2, duration=$3, photo_id=$4, icon_id=$5, updated_at=$6,price=$7 where id=$8`
	_, err = tx.Exec(statement, model.Name, model.Description, model.Duration, model.PhotoID.String, model.IconID.String, model.UpdatedAt, model.Price.Int32, model.ID)
	if err != nil {
		return err
	}

	return nil
}

func (TreatmentRepository) Add(model models.Treatment, tx *sql.Tx) (res string, err error) {
	statement := `insert into "master_treatments" (name,description,duration,price,photo_id,icon_id,created_at,updated_at,category_id) values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`
	err = tx.QueryRow(statement, model.Name, model.Description, model.Duration, model.Price.Int32, model.PhotoID.String, model.IconID.String, model.CreatedAt, model.UpdatedAt,model.CategoryID).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (TreatmentRepository) Delete(model models.Treatment, tx *sql.Tx) (err error) {
	statement := `update "master_treatments" set updated_at=$1, deleted_at=$2 where id=$3`
	_, err = tx.Exec(statement, model.UpdatedAt, model.DeletedAt.Time, model.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repository TreatmentRepository) CountBy(ID, column, value,categoryId, operator string) (res int, err error) {
	whereStatement := `where ` + column + `` + operator + `$1 and t.deleted_at is null and t.category_id=$2`
	whereParams := []interface{}{value,categoryId}
	if ID != "" {
		whereStatement = `where (` + column + `` + operator + `$1 and t.deleted_at is null and t.category_id=$2) and t.id <>$3`
		whereParams = append(whereParams, ID)
	}
	statement := `select count(distinct t.id) from "master_treatments" t ` + treatmentJoinStatement + ` ` + whereStatement
	err = repository.DB.QueryRow(statement, whereParams...).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}
