package actions

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/db/repositories/contracts"
	"strings"
)

type BedRepository struct {
	DB *sql.DB
}

func NewBedRepository(DB *sql.DB) contracts.IBedRepository {
	bedWhereParams = []interface{}{}
	return &BedRepository{DB: DB}
}

const (
	bedSelectStatement = `select id,bed_code,clinic_id,name,address,pic_name,phone_number,email,is_use_able,treatments,created_at,updated_at,deleted_at`
)

var (
	bedWhereParams    = []interface{}{}
	bedWhereStatement = `where (bed_code like $1 or lower(name) like $1) and deleted_at is null`
)

func (repository BedRepository) scanRows(rows *sql.Rows) (res models.Bed, err error) {
	err = rows.Scan(&res.ID, &res.BedCode, &res.Clinic.ID, &res.Clinic.Name, &res.Clinic.Address, &res.Clinic.PICName, &res.Clinic.PhoneNumber, &res.Clinic.Email, &res.IsUseAble, &res.Treatments, &res.CreatedAt, &res.UpdatedAt,
		&res.DeletedAt)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository BedRepository) scanRow(row *sql.Row) (res models.Bed, err error) {
	err = row.Scan(&res.ID, &res.BedCode, &res.Clinic.ID, &res.Clinic.Name, &res.Clinic.Address, &res.Clinic.PICName, &res.Clinic.PhoneNumber, &res.Clinic.Email, &res.IsUseAble, &res.Treatments, &res.CreatedAt, &res.UpdatedAt,
		&res.DeletedAt)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository BedRepository) Browse(clinicID, search, order, sort string, limit, offset int) (data []models.Bed, count int, err error) {
	bedWhereParams = []interface{}{"%" + strings.ToLower(search) + "%"}
	if clinicID != "" {
		bedWhereStatement += ` and clinic_id=$2`
		bedWhereParams = append(bedWhereParams, clinicID)
	}
	bedWhereParams = append(bedWhereParams, []interface{}{limit, offset}...)
	statement := bedSelectStatement + ` from "clinic_beds" ` + bedWhereStatement + ` order by ` + order + ` ` + sort + ` limit $3 offset $4`
	rows, err := repository.DB.Query(statement, bedWhereParams...)
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

	statement = `select count(id) from "clinic_beds" ` + bedWhereStatement
	err = repository.DB.QueryRow(statement, bedWhereParams[0], bedWhereParams[1]).Scan(&count)
	if err != nil {
		return data, count, err
	}

	return data, count, nil
}

func (repository BedRepository) BrowseAll(clinicID, search string) (data []models.Bed, err error) {
	bedWhereParams = []interface{}{"%" + strings.ToLower(search) + "%"}
	if clinicID != "" {
		bedWhereStatement += ` and clinic_id=$2`
		bedWhereParams = append(bedWhereParams, clinicID)
	}
	statement := bedSelectStatement + ` from "clinic_beds" ` + bedWhereStatement
	rows, err := repository.DB.Query(statement, bedWhereParams...)
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

func (repository BedRepository) ReadBy(column, value, operator string) (data models.Bed, err error) {
	statement := bedSelectStatement + ` from "clinic_beds" where ` + column + `` + operator + `$1 and "deleted_at" is null`
	row := repository.DB.QueryRow(statement, value)
	data, err = repository.scanRow(row)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (BedRepository) Edit(input models.Bed, tx *sql.Tx) (err error) {
	statement := `update "beds" set bed_code=$1, clinic_id=$2, is_use_able=$3, updated_at=$4 where id=$5 returning id`
	_, err = tx.Exec(statement, input.BedCode, input.ClinicID, input.IsUseAble, input.UpdatedAt, input.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repository BedRepository) Add(input models.Bed, tx *sql.Tx) (res string, err error) {
	statement := `insert into "beds" (bed_code,clinic_id,is_use_able,created_at,updated_at) values($1,$2,$3,$4,$5) returning id`
	err = repository.DB.QueryRow(statement, input.BedCode, input.ClinicID, input.IsUseAble, input.CreatedAt, input.UpdatedAt).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (BedRepository) DeleteBy(column, value, operator string, input models.Bed, tx *sql.Tx) (err error) {
	statement := `update "beds" set updated_at=$1, deleted_at=$2 where ` + column + `` + operator + `$3 returning id`
	_, err = tx.Exec(statement, input.UpdatedAt, input.DeletedAt.Time, value)
	if err != nil {
		return err
	}

	return nil
}

func (repository BedRepository) CountBy(ID, column, value, operator string) (res int, err error) {
	countWhereStatement := `where ` + column + `` + operator + `$1 and deleted_at is null`
	countWhereParams := []interface{}{value}
	if ID != "" {
		countWhereStatement = `where (` + column + `` + operator + `$1 and deleted_at is null) and id <> $2`
		countWhereParams = []interface{}{value, ID}
	}
	statement := `select count(id) from "clinic_beds" ` + countWhereStatement
	err = repository.DB.QueryRow(statement, countWhereParams...).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}
