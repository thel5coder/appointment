package actions

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/db/repositories/contracts"
	"strings"
)

type ProcedureRepository struct {
	DB *sql.DB
}

func NewProcedureRepository(DB *sql.DB) contracts.IProcedureRepository {
	return &ProcedureRepository{DB: DB}
}

const (
	procedureSelectStatement = `select id,name,duration,created_at,updated_at,deleted_at`
	procedureWhereStatement  = `where (lower(name) like $1) and deleted_at is null`
)

func (repository ProcedureRepository) scanRows(rows *sql.Rows) (res models.Procedure, err error) {
	err = rows.Scan(&res.ID, &res.Name, &res.Duration, &res.CreatedAt, &res.UpdatedAt, &res.DeletedAt)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository ProcedureRepository) scanRow(row *sql.Row) (res models.Procedure, err error) {
	err = row.Scan(&res.ID, &res.Name, &res.Duration, &res.CreatedAt, &res.UpdatedAt, &res.DeletedAt)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository ProcedureRepository) Browse(search, order, sort string, limit, offset int) (data []models.Procedure, count int, err error) {
	statement := procedureSelectStatement + ` from "procedures" ` + procedureWhereStatement + ` order by ` + order + ` ` + sort + ` limit $2 offset $3`

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

	statement = `select count(id) from "procedures" ` + procedureWhereStatement
	err = repository.DB.QueryRow(statement, "%"+strings.ToLower(search)+"%").Scan(&count)
	if err != nil {
		return data, count, err
	}

	return data, count, nil
}

func (repository ProcedureRepository) BrowseAll(search string) (data []models.Procedure, err error) {
	statement := procedureSelectStatement + ` from "procedures" ` + procedureWhereStatement

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

	return data, nil
}

func (repository ProcedureRepository) ReadBy(column, value, operator string) (data models.Procedure, err error) {
	statement := procedureSelectStatement + ` from "procedures" where ` + column + `` + operator + `$1 and "deleted_at" is null`

	row := repository.DB.QueryRow(statement, value)
	data, err = repository.scanRow(row)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (repository ProcedureRepository) Edit(model models.Procedure) (res string, err error) {
	statement := `update "procedures" set "name"=$1, duration=$2, updated_at=$3 where id=$4 returning id`

	err = repository.DB.QueryRow(statement, model.Name, model.Duration, model.UpdatedAt, model.ID).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository ProcedureRepository) Add(model models.Procedure) (res string, err error) {
	statement := `insert into "procedures" (name,duration,created_at,updated_at) values($1,$2,$3,$4) returning id`

	err = repository.DB.QueryRow(statement, model.Name, model.Duration, model.CreatedAt, model.UpdatedAt).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository ProcedureRepository) DeleteBy(column, value, operator string, model models.Procedure) (res string, err error) {
	statement := `update "procedures" set updated_at=$1, deleted_at=$2 where ` + column + `` + operator + `$3 returning id`

	err = repository.DB.QueryRow(statement, model.UpdatedAt, model.DeletedAt.Time, value).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository ProcedureRepository) CountBy(ID, column, value, operator string) (res int, err error) {
	conditionWhereStatement := `where ` + column + `` + operator + `$1 and deleted_at is null`
	conditionWhereParams := []interface{}{value}
	if ID != "" {
		conditionWhereStatement = `where (` + column + `` + operator + `$1 and deleted_at is null) and id <>$2`
		conditionWhereParams = append(conditionWhereParams, ID)
	}
	statement := `select count(id) from "procedures" ` + conditionWhereStatement

	err = repository.DB.QueryRow(statement, conditionWhereParams...).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}
