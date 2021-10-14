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

type ClinicRepository struct {
	DB *sql.DB
}

var (
	clinicSelectStatementParams = []interface{}{}
	clinicSetUpdateStatement    = ``
	clinicUpdateParams          = []interface{}{}
)

const (
	clinicSelect         = `select "id","name","address","pic_name","phone_number","email","created_at","updated_at"`
	clinicWhereStatement = `where(lower("name") like $1 or lower("address") like $1 or lower("pic_name") like $1 or lower("email") like $1 or "phone_number" like $1) 
                           and "deleted_at" is null`
)

func NewClinicRepository(DB *sql.DB) contracts.IClinicRepository {
	clinicSelectStatementParams = []interface{}{}
	clinicSetUpdateStatement = ``
	clinicUpdateParams = []interface{}{}

	return &ClinicRepository{DB: DB}
}

func (repository ClinicRepository) scanRows(rows *sql.Rows) (res models.Clinic, err error) {
	err = rows.Scan(&res.ID, &res.Name, &res.Address, &res.PICName, &res.PhoneNumber, &res.Email, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository ClinicRepository) scanRow(row *sql.Row) (res models.Clinic, err error) {
	err = row.Scan(&res.ID, &res.Name, &res.Address, &res.PICName, &res.PhoneNumber, &res.Email, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository ClinicRepository) Browse(search, order, sort string, limit, offset int) (data []models.Clinic, count int, err error) {
	clinicSelectStatementParams = append(clinicSelectStatementParams, []interface{}{"%" + strings.ToLower(search) + "%", limit, offset}...)
	statement := clinicSelect + ` from "master_clinics" ` + clinicWhereStatement + ` order by ` + order + ` ` + sort + ` limit $2 offset $3`
	rows, err := repository.DB.Query(statement, clinicSelectStatementParams...)
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

	statement = `select count("id") from "master_clinics" ` + clinicWhereStatement
	err = repository.DB.QueryRow(statement, "%"+strings.ToLower(search)+"%").Scan(&count)
	if err != nil {
		return data, count, err
	}
	return data, count, err
}

func (repository ClinicRepository) BrowseAll(search string) (data []models.Clinic, err error) {
	statement := clinicSelect + ` from "master_clinics" ` + clinicWhereStatement
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

func (repository ClinicRepository) ReadBy(column, value, operator string) (data models.Clinic, err error) {
	statement := clinicSelect + ` from "master_clinics" where ` + column + `` + operator + `$1 and "deleted_at" is null`
	row := repository.DB.QueryRow(statement, value)
	data, err = repository.scanRow(row)
	if err != nil {
		return data, err
	}

	return data, err
}

func (repository ClinicRepository) Edit(input viewmodel.ClinicVm) (res string, err error) {
	statement := `update "master_clinics" set "name"=$1,"address"=$2, "pic_name"=$3, "phone_number"=$4, "email"=$5, "updated_at"=$6 where "id"=$7 returning "id"`
	err = repository.DB.QueryRow(statement, input.Name, input.Address, input.PICName, input.PhoneNumber, input.Email,
		datetime.StrParseToTime(input.UpdatedAt, time.RFC3339), &input.ID).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository ClinicRepository) Add(input viewmodel.ClinicVm) (res string, err error) {
	statement := `insert into "master_clinics" ("name","address","pic_name","phone_number","email","created_at","updated_at") values($1,$2,$3,$4,$5,$6,$7) returning "id"`
	err = repository.DB.QueryRow(statement, input.Name, input.Address, input.PICName, input.PhoneNumber, input.Email, datetime.StrParseToTime(input.CreatedAt, time.RFC3339),
		datetime.StrParseToTime(input.UpdatedAt, time.RFC3339)).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository ClinicRepository) Delete(ID, updatedAt, deletedAt string) (res string, err error) {
	statement := `update "master_clinics" set "updated_at"=$1, "deleted_at"=$2 where "id"=$3 returning "id"`
	err = repository.DB.QueryRow(statement, datetime.StrParseToTime(updatedAt, time.RFC3339), datetime.StrParseToTime(deletedAt, time.RFC3339), ID).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository ClinicRepository) CountBy(ID, column, value, operator string) (res int, err error) {
	clinicSelectStatementParams = append(clinicSelectStatementParams, value)
	whereStatement := `where ` + column + `` + operator + `$1 and "deleted_at" is null`
	if ID != "" {
		clinicSelectStatementParams = append(clinicSelectStatementParams, ID)
		whereStatement = `where (` + column + `` + operator + `$1 and "deleted_at" is null) and "id"<>$2`
	}
	statement := `select count("id") from "master_clinics" ` + whereStatement
	err = repository.DB.QueryRow(statement, clinicSelectStatementParams...).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}
