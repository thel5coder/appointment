package actions

import (
	"database/sql"
	"fmt"
	"profira-backend/db/models"
	"profira-backend/db/repositories/contracts"
	"profira-backend/helpers/datetime"
	"profira-backend/usecase/viewmodel"
	"strings"
	"time"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(DB *sql.DB) contracts.IUserRepository {
	userSelectStatementParams = []interface{}{}
	userSetUpdateStatement = ``
	userUpdateParams = []interface{}{}

	return &UserRepository{DB: DB}
}

const (
	userSelect = `select u."id",u."name",u."email", u."mobile_phone",u."password",u."fcm_device_token",u."is_active",r."id",r."slug",r."name",u."profile_picture_id",f."path",u."created_at",u."updated_at",
                  u."activated_at"
                  from "users" u`
	userInnerJoin = `inner join "roles" r on r."id"=u."role_id" and r."deleted_at" is null
                    left join "files" f on f."id"=u."profile_picture_id" and f."deleted_at" is null`
	userWhereStatement = `where (lower(u."name") like $1 or lower(u."email") like $1 or lower(u."mobile_phone") like $1) and u."deleted_at" is null`
)

var (
	userSelectStatementParams = []interface{}{}
	userSetUpdateStatement    = ``
	userUpdateParams          = []interface{}{}
)

func (UserRepository) scanRows(rows *sql.Rows) (res models.User, err error) {
	err = rows.Scan(
		&res.ID,
		&res.Name,
		&res.Email,
		&res.MobilePhone,
		&res.Password,
		&res.FcmDeviceToken,
		&res.IsActive,
		&res.Role.ID,
		&res.Role.Slug,
		&res.Role.Name,
		&res.ProfilePictureID,
		&res.ProfilePicturePath,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.ActivatedAt,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (UserRepository) scanRow(row *sql.Row) (res models.User, err error) {
	err = row.Scan(
		&res.ID,
		&res.Name,
		&res.Email,
		&res.MobilePhone,
		&res.Password,
		&res.FcmDeviceToken,
		&res.IsActive,
		&res.Role.ID,
		&res.Role.Slug,
		&res.Role.Name,
		&res.ProfilePictureID,
		&res.ProfilePicturePath,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.ActivatedAt,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository UserRepository) Browse(search, order, sort string, limit, offset int) (data []models.User, count int, err error) {
	userSelectStatementParams = append(userSelectStatementParams, []interface{}{"%" + strings.ToLower(search) + "%", limit, offset}...)

	statement := userSelect + ` ` + userInnerJoin + ` ` + userWhereStatement + ` order by u.` + order + ` ` + sort + ` limit $2 offset $3`
	rows, err := repository.DB.Query(statement, userSelectStatementParams...)
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

	statement = `select count(u."id") from "users" u ` + userInnerJoin + ` ` + userWhereStatement
	err = repository.DB.QueryRow(statement, "%"+strings.ToLower(search)+"%").Scan(&count)
	if err != nil {
		return data, count, err
	}

	return data, count, nil
}

func (repository UserRepository) BrowseAllBy(column, value string) (data []models.User, err error) {
	whereStatement := `u."deleted_at" is null`
	if column != "" {
		whereStatement = `where ` + column + `=$1 and u."deleted_at" is null`
		userSelectStatementParams = append(userSelectStatementParams, value)
	}
	statement := userSelect + ` ` + userInnerJoin + ` ` + whereStatement
	rows, err := repository.DB.Query(statement, userSelectStatementParams...)
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

func (repository UserRepository) ReadBy(column, value, operator string) (data models.User, err error) {
	statement := userSelect + ` ` + userInnerJoin + ` where ` + column + `` + operator + `$1 and u."deleted_at" is null`
	row := repository.DB.QueryRow(statement, value)
	data, err = repository.scanRow(row)
	if err != nil {
		fmt.Println(err.Error())
		return data, err
	}

	return data, nil
}

func (UserRepository) Edit(input viewmodel.UserVm, password string, tx *sql.Tx) (err error) {
	userUpdateParams = append(userUpdateParams, []interface{}{
		input.Name,
		input.Email,
		input.MobilePhone,
		input.ProfilePicture.ID,
		input.IsActive,
		datetime.StrParseToTime(input.UpdatedAt, time.RFC3339),
		input.ID,
	}...)
	userSetUpdateStatement = `"name"=$1, "email"=$2, "mobile_phone"=$3, "profile_picture_id"=$4, "is_active"=$5, "updated_at"=$6 where "id"=$7`
	if password != "" {
		userSetUpdateStatement = `"name"=$1, "email"=$2, "mobile_phone"=$3, "profile_picture_id"=$4, "is_active"=$5,"password"=$8, "updated_at"=$6 where "id"=$7`
		userUpdateParams = append(userUpdateParams, []interface{}{password}...)
	}
	statement := `update "users" set ` + userSetUpdateStatement
	_, err = tx.Exec(statement, userUpdateParams...)
	if err != nil {
		return err
	}

	return nil
}

func (repository UserRepository) EditPassword(ID, password, updatedAt string) (res string, err error) {
	statement := `update "users" set "password"=$1, "updated_at"=$2 where "id"=$3 returning "id"`
	err = repository.DB.QueryRow(statement, password, datetime.StrParseToTime(updatedAt, time.RFC3339), ID).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository UserRepository) EditFcmDeviceToken(ID, fcmDeviceToken, updatedAt string) (res string, err error) {
	statement := `update "users" set "fcm_device_token"=$1, "updated_at"=$2 where "id"=$3 returning "id"`
	err = repository.DB.QueryRow(statement, fcmDeviceToken, datetime.StrParseToTime(updatedAt, time.RFC3339), ID).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository UserRepository) EditActivatedUser(ID, updatedAt, activatedAt string, isActive bool) (res string, err error) {
	statement := `update "users" set "is_active"=$1,"updated_at"=$2, "activated_at"=$3 where "id"=$4 returning "id"`
	err = repository.DB.QueryRow(statement, isActive, datetime.StrParseToTime(updatedAt, time.RFC3339), datetime.StrParseToTime(activatedAt, time.RFC3339), ID).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (UserRepository) Add(input viewmodel.UserVm, password string, tx *sql.Tx) (res string, err error) {
	statement := `insert into "users" (name,"email","mobile_phone","password","is_active",profile_picture_id,"role_id","activated_at","created_at","updated_at") 
                  values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) returning "id"`
	err = tx.QueryRow(
		statement,
		input.Name,
		input.Email,
		input.MobilePhone,
		password,
		input.IsActive,
		input.ProfilePicture.ID,
		input.Role.ID,
		datetime.EmptyTime(datetime.StrParseToTime(input.ActivatedAt, time.RFC3339)),
		datetime.StrParseToTime(input.CreatedAt, time.RFC3339),
		datetime.StrParseToTime(input.UpdatedAt, time.RFC3339),
	).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (UserRepository) DeleteBy(column, value, operator, updatedAt, deletedAt string, tx *sql.Tx) (err error) {
	statement := `update "users" set "updated_at"=$1,"deleted_at"=$2 where ` + column + `` + operator + `$3`
	_, err = tx.Exec(statement, datetime.StrParseToTime(updatedAt, time.RFC3339), datetime.StrParseToTime(deletedAt, time.RFC3339), value)
	if err != nil {
		return err
	}

	return nil
}

func (repository UserRepository) CountBy(ID, column, value, operator string) (res int, err error) {
	userSelectStatementParams = append(userSelectStatementParams, value)
	whereStatement := `where ` + column + `` + operator + `$1 and "deleted_at" is null`
	if ID != "" {
		whereStatement = `where (` + column + `` + operator + `$1 and "deleted_at" is null) and "id"<>$2`
		userSelectStatementParams = append(userSelectStatementParams, ID)
	}

	statement := `select count("id") from "users" ` + whereStatement
	err = repository.DB.QueryRow(statement, userSelectStatementParams...).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}
