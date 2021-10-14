package actions

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/db/repositories/contracts"
	"profira-backend/helpers/datetime"
	"profira-backend/usecase/viewmodel"
	"strings"
	"time"
)

type RoleRepository struct {
	DB *sql.DB
}

func NewRoleRepository(DB *sql.DB) contracts.IRoleRepository {
	roleSelectStatementParams = []interface{}{}

	return &RoleRepository{DB: DB}
}

const (
	roleSelect         = `select "id","slug","name","created_at","updated_at","deleted_at"`
	roleWhereStatement = `where lower("name") like $1 and "deleted_at" is null`
)

var (
	roleSelectStatementParams = []interface{}{}
)

func (repository RoleRepository) scanRow(row *sql.Row) (res models.Role, err error) {
	err = row.Scan(
		&res.ID,
		&res.Slug,
		&res.Name,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.DeletedAt,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository RoleRepository) scanRows(rows *sql.Rows) (res models.Role, err error) {
	err = rows.Scan(
		&res.ID,
		&res.Slug,
		&res.Name,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.DeletedAt,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository RoleRepository) Browse(search, order, sort string, limit, offset int) (data []models.Role, count int, err error) {
	roleSelectStatementParams = append(roleSelectStatementParams, []interface{}{"%" + strings.ToLower(search) + "%", limit, offset}...)

	statement := roleSelect + ` from "roles" ` + roleWhereStatement + ` order by ` + order + ` ` + sort + ` limit $2 offset $3`
	rows, err := repository.DB.Query(statement, roleSelectStatementParams...)
	if err != nil {
		return data, count, err
	}

	for rows.Next() {
		dataTemp, err := repository.scanRows(rows)
		if err != nil {
			return data, count, err
		}
		data = append(data, dataTemp)
	}

	statement = `select count("id") from "roles" ` + roleWhereStatement
	err = repository.DB.QueryRow(statement, "%"+strings.ToLower(search)+"%").Scan(&count)
	if err != nil {
		return data, count, err
	}

	return data, count, nil
}

func (repository RoleRepository) ReadBy(column, value string) (data models.Role, err error) {
	statement := roleSelect + `from "roles"
                 where ` + column + `=$1 and "deleted_at" is null`
	row := repository.DB.QueryRow(statement, value)
	if data, err = repository.scanRow(row); err != nil {
		return data, err
	}

	return data, nil
}

func (repository RoleRepository) Edit(input viewmodel.RoleVm) (res string, err error) {
	statement := `update "roles" set "name"=$1, "slug"=$2, "updated_at"=$3 where "id"=$4 returning "id"`
	err = repository.DB.QueryRow(statement, input.Name, input.Slug, datetime.StrParseToTime(input.UpdatedAt, time.RFC3339), input.ID).Scan(&res)

	return res, err
}

func (repository RoleRepository) Add(input viewmodel.RoleVm) (res string, err error) {
	statement := `insert into "roles" ("name","slug","created_at","updated_at") values($1,$2,$3,$4) returning "id"`
	err = repository.DB.QueryRow(statement, input.Name, input.Slug, datetime.StrParseToTime(input.CreatedAt, time.RFC3339), datetime.StrParseToTime(input.UpdatedAt, time.RFC3339)).Scan(&res)

	return res, err
}

func (repository RoleRepository) Delete(ID, updatedAt, deletedAt string) (res string, err error) {
	statement := `update "roles" set "updated_at"=$1, "deleted_at"=$2 where "id"=$3 returning "id"`
	err = repository.DB.QueryRow(statement, datetime.StrParseToTime(updatedAt, time.RFC3339), datetime.StrParseToTime(deletedAt, time.RFC3339), ID).Scan(&res)

	return res, err
}

func (repository RoleRepository) CountBy(ID, column, value string) (res int, err error) {
	if ID == "" {
		statement := `select count("id") from "roles" where ` + column + `=$1 and "deleted_at" is null`
		err = repository.DB.QueryRow(statement, value).Scan(&res)
	} else {
		statement := `select count("id") from "roles" where (` + column + `=$1 and "deleted_at" is null)and "id"<>$2`
		err = repository.DB.QueryRow(statement, value, ID).Scan(&res)
	}

	return res, err
}
