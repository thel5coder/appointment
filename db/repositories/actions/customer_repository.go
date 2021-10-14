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

type CustomerRepository struct {
	DB *sql.DB
}

func NewCustomerRepository(DB *sql.DB) contracts.ICustomerRepository {
	customerSelectStatementParams = []interface{}{}
	customerSetUpdateStatement = ``
	customerUpdateParams = []interface{}{}

	return &CustomerRepository{DB: DB}
}

const (
	customerSelect = `select c."id",c."name",c."sex",c."address",c."birth_date",c."marital_status",c."phone_number",c."mobile_phone_1",
                      c."mobile_phone_2",c."religion",c."education",c."hobby",c."profession",c."reference",c."notes",c."city_id",c."user_id",
                      c."created_at",c."updated_at",c."deleted_at",u."id",u."email",u."mobile_phone",u."role_id",u."profile_picture_id",
                      u."is_active",f."path"`
	customerInnerJoin = `inner join "users" u on u."id"=c."user_id" and u."deleted_at" is null
                         left join "files" f on f."id"=u."profile_picture_id" and f."deleted_at" is null`
	customerWhereStatement = `where (lower(c."name") like $1 or cast(c."sex" as varchar) like $1 or u."mobile_phone" like $1 or lower(u."email") like $1) and c."deleted_at" is null`
)

var (
	customerSelectStatementParams = []interface{}{}
	customerSetUpdateStatement    = ``
	customerUpdateParams          = []interface{}{}
)

func (repository CustomerRepository) scanRows(rows *sql.Rows) (res models.Customer, err error) {
	err = rows.Scan(&res.ID, &res.Name, &res.Sex, &res.Address, &res.BirthDate, &res.MaritalStatus, &res.PhoneNumber, &res.MobilePhoneNumber1, &res.MobilePhoneNumber2, &res.Religion, &res.Education,
		&res.Hobby, &res.Profession, &res.Reference, &res.Notes, &res.CityID, &res.UserID, &res.CreatedAt, &res.UpdatedAt, &res.DeletedAt, &res.User.ID, &res.User.Email,
		&res.User.MobilePhone, &res.User.Role.ID, &res.User.ProfilePictureID, &res.User.IsActive, &res.User.ProfilePicturePath)
	if err != nil {
		return res, err
	}

	return res, err
}

func (repository CustomerRepository) scanRow(row *sql.Row) (res models.Customer, err error) {
	err = row.Scan(&res.ID, &res.Name, &res.Sex, &res.Address, &res.BirthDate, &res.MaritalStatus, &res.PhoneNumber, &res.MobilePhoneNumber1, &res.MobilePhoneNumber2, &res.Religion, &res.Education,
		&res.Hobby, &res.Profession, &res.Reference, &res.Notes, &res.CityID, &res.UserID, &res.CreatedAt, &res.UpdatedAt, &res.DeletedAt, &res.User.ID, &res.User.Email,
		&res.User.MobilePhone, &res.User.Role.ID, &res.User.ProfilePictureID, &res.User.IsActive, &res.User.ProfilePicturePath)
	if err != nil {
		return res, err
	}

	return res, err
}

func (repository CustomerRepository) Browse(search, order, sort string, limit, offset int) (data []models.Customer, count int, err error) {
	customerSelectStatementParams = append(customerSelectStatementParams, []interface{}{"%" + strings.ToLower(search) + "%", limit, offset}...)
	statement := customerSelect + ` from "customers" c ` + customerInnerJoin + ` ` + customerWhereStatement + ` order by c.` + order + ` ` + sort + ` limit $2 offset $3`
	rows, err := repository.DB.Query(statement, customerSelectStatementParams...)
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

	statement = `select count(c."id") from "customers" c ` + customerInnerJoin + ` ` + customerWhereStatement
	err = repository.DB.QueryRow(statement, "%"+strings.ToLower(search)+"%").Scan(&count)
	if err != nil {
		return data, count, err
	}

	return data, count, nil
}

func (repository CustomerRepository) BrowseAllBy(column, value string) (data []models.Customer, err error) {
	whereStatement := `c."deleted_at" is null`
	if column != "" {
		whereStatement = `where ` + column + `=$1 and c."deleted_at" is null`
		customerSelectStatementParams = append(customerSelectStatementParams, value)
	}
	statement := customerSelect + ` from "customers" c ` + customerInnerJoin + ` ` + whereStatement
	rows, err := repository.DB.Query(statement, customerSelectStatementParams...)
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

func (repository CustomerRepository) ReadBy(column, value, operator string) (data models.Customer, err error) {
	statement := customerSelect + ` from "customers" c ` + customerInnerJoin + ` where ` + column + `` + operator + `$1`
	row := repository.DB.QueryRow(statement, value)
	data, err = repository.scanRow(row)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (CustomerRepository) Edit(input viewmodel.CustomerVm, tx *sql.Tx) (err error) {
	statement := `update "customers" set "name"=$1, "birth_date"=$2,"updated_at"=$3, "sex"=$4, address=$5 where "id"=$6`
	_, err = tx.Exec(statement, input.Name, input.BirthDate, datetime.StrParseToTime(input.BirthDate, "2006-01-02"), input.Sex, input.Address, input.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repository CustomerRepository) EditProfile(model models.Customer, tx *sql.Tx) (err error) {
	statement := `update customers set name=$1, birth_date=$2, updated_at=$3 where user_id=$4`
	_, err = tx.Exec(statement, model.Name, model.BirthDate, model.UpdatedAt, model.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (CustomerRepository) Add(input viewmodel.CustomerVm, tx *sql.Tx) (err error) {
	statement := `insert into "customers" ("user_id","name","sex","address","birth_date","created_at","updated_at")
                 values($1,$2,$3,$4,$5,$6,$7)`
	_, err = tx.Exec(
		statement,
		input.User.ID,
		input.Name,
		input.Sex,
		input.Address,
		datetime.StrParseToTime(input.BirthDate, "2006-01-2"),
		datetime.StrParseToTime(input.CreatedAt, time.RFC3339),
		datetime.StrParseToTime(input.UpdatedAt, time.RFC3339),
	)
	if err != nil {
		return err
	}

	return nil
}

func (CustomerRepository) DeleteBy(column, value, operator, updatedAt, deletedAt string, tx *sql.Tx) (err error) {
	statement := `update "customers" set "updated_at"=$1, "deleted_at"=$2 where ` + column + `` + operator + `$3`
	_, err = tx.Exec(statement, datetime.StrParseToTime(updatedAt, time.RFC3339), datetime.StrParseToTime(deletedAt, time.RFC3339), value)
	if err != nil {
		return err
	}

	return nil
}

func (repository CustomerRepository) CountBy(ID, column, value, operator string) (res int, err error) {
	customerSelectStatementParams = append(customerSelectStatementParams, value)
	whereStatement := `where ` + column + `` + operator + `$1 and "deleted_at" is null`
	if ID != "" {
		whereStatement = `where (` + column + `` + operator + `$1 and "deleted_at" is null) and "id"<>$2`
		customerSelectStatementParams = append(customerSelectStatementParams, ID)
	}
	statement := `select count("id") from "customers" ` + whereStatement
	err = repository.DB.QueryRow(statement, customerSelectStatementParams...).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}
