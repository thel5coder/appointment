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

type StaffRepository struct {
	DB *sql.DB
}

func NewStaffRepository(DB *sql.DB) contracts.IStaffRepository {
	staffSelectStatementParams = []interface{}{}
	staffSetUpdateStatement = ``
	staffUpdateParams = []interface{}{}

	return &StaffRepository{DB: DB}
}

const (
	staffSelect = `select s."id",s."name",s."description",s."user_id",s."created_at",s."updated_at",u."id",u."email",u."mobile_phone",
                   u."role_id",u."profile_picture_id",u."is_active",f."path",
                  array_to_string(array_agg(mc."id" || ':' || mc."name" || ':' || mc."address" || ':' || mc."pic_name" || ':' || mc."phone_number" || ':' || mc."email"),','),
                  array_to_string(array_agg(dt.id || ':' || dt.treatment_id || ':' || mt.name),',')`
	staffJoin = `inner join "users" u on u."id"=s."user_id" and u."deleted_at" is null
                 left join "files" f on f."id"=u."profile_picture_id"
                 inner join "staff_clinics" sc on sc."staff_id"=s."id"
                 inner join "master_clinics" mc on mc."id"=sc."clinic_id"
                 left join "doctor_treatments" dt on dt.doctor_id=s.id and dt.deleted_at is null
                 left join "master_treatments" mt on mt.id=dt.treatment_id and mt.deleted_at is null`
	staffGroupBy        = `group by s."id",u."id",f."path"`
	staffWhereStatement = `where (lower(s."name") like $1 or lower(u."email") like $1 or lower(mc."name") like $1 or u."mobile_phone" like $1) 
                          and s."deleted_at" is null and u."role_id"=$2`
)

var (
	staffSelectStatementParams = []interface{}{}
	staffSetUpdateStatement    = ``
	staffUpdateParams          = []interface{}{}
)

func (repository StaffRepository) scanRows(rows *sql.Rows) (data models.Staff, err error) {
	err = rows.Scan(&data.ID, &data.Name, &data.Description, &data.UserID, &data.CreatedAt, &data.UpdatedAt, &data.User.ID, &data.User.Email, &data.User.MobilePhone, &data.User.Role.ID,
		&data.User.ProfilePictureID, &data.User.IsActive, &data.User.ProfilePicturePath, &data.Clinics, &data.Treatments)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (repository StaffRepository) scanRow(row *sql.Row) (data models.Staff, err error) {
	err = row.Scan(&data.ID, &data.Name, &data.Description, &data.UserID, &data.CreatedAt, &data.UpdatedAt, &data.User.ID, &data.User.Email, &data.User.MobilePhone, &data.User.Role.ID,
		&data.User.ProfilePictureID, &data.User.IsActive, &data.User.ProfilePicturePath, &data.Clinics, &data.Treatments)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (repository StaffRepository) BrowseByRole(roleID, search, order, sort string, limit, offset int) (data []models.Staff, count int, err error) {
	staffSelectStatementParams = append(staffSelectStatementParams, []interface{}{"%" + strings.ToLower(search) + "%", roleID, limit, offset}...)
	statement := staffSelect + ` from "master_staffs" s ` + staffJoin + ` ` + staffWhereStatement + ` ` + staffGroupBy + ` order by s.` + order + ` ` + sort + ` limit $3 offset $4`
	rows, err := repository.DB.Query(statement, staffSelectStatementParams...)
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

	statement = `select distinct count(s."id") from "master_staffs" s ` + staffJoin + ` ` + staffWhereStatement
	err = repository.DB.QueryRow(statement, "%"+strings.ToLower(search)+"%", roleID).Scan(&count)
	if err != nil {
		return data, count, err
	}

	return data, count, nil
}

func (repository StaffRepository) BrowseStaffDoctorByClinic(clinicId,roleId string) (data []models.Staff, err error) {
	statement := staffSelect + ` FROM master_staffs s ` + staffJoin + ` WHERE sc.clinic_id =$1 and u.role_id=$2` +  staffGroupBy + ` ORDER BY s.name ASC`
	rows, err := repository.DB.Query(statement, clinicId,roleId)
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

func (repository StaffRepository) BrowseAllByRole(roleID, search string) (data []models.Staff, err error) {
	staffSelectStatementParams = append(staffSelectStatementParams, []interface{}{"%" + strings.ToLower(search) + "%", roleID}...)
	statement := staffSelect + ` from "master_staffs" s ` + staffJoin + ` ` + staffWhereStatement + ` ` + staffGroupBy
	rows, err := repository.DB.Query(statement, staffSelectStatementParams...)
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

func (repository StaffRepository) ReadBy(column, value, operator string) (data models.Staff, err error) {
	statement := staffSelect + ` from "master_staffs" s ` + staffJoin + ` where ` + column + `` + operator + `$1 and s."deleted_at" is null ` + staffGroupBy
	row := repository.DB.QueryRow(statement, value)
	data, err = repository.scanRow(row)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (StaffRepository) Edit(input viewmodel.StaffVm, tx *sql.Tx) (err error) {
	statement := `update "master_staffs" set "name"=$1, "description"=$2, "updated_at"=$3 where "id"=$4`
	_, err = tx.Exec(statement, input.Name, input.Description, datetime.StrParseToTime(input.UpdatedAt, time.RFC3339), input.ID)
	if err != nil {
		return err
	}

	return nil
}

func (StaffRepository) Add(input viewmodel.StaffVm, tx *sql.Tx) (res string, err error) {
	statement := `insert into "master_staffs" ("name","description","user_id","created_at","updated_at") values($1,$2,$3,$4,$5) returning "id"`
	err = tx.QueryRow(statement, input.Name, input.Description, input.User.ID, datetime.StrParseToTime(input.CreatedAt, time.RFC3339),
		datetime.StrParseToTime(input.UpdatedAt, time.RFC3339)).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (StaffRepository) DeleteBy(column, value, operator, updatedAt, deletedAt string, tx *sql.Tx) (err error) {
	statement := `update "master_staffs" set "updated_at"=$1, "deleted_at"=$2 where ` + column + `` + operator + `$3`
	_, err = tx.Exec(statement, datetime.StrParseToTime(updatedAt, time.RFC3339), datetime.StrParseToTime(deletedAt, time.RFC3339), value)
	if err != nil {
		return err
	}

	return nil
}

func (repository StaffRepository) CountBy(ID, column, value, operator string) (res int, err error) {
	staffSelectStatementParams = append(staffSelectStatementParams, value)
	whereStatement := `where ` + column + `` + operator + `$1 and s."deleted_at" is null`
	if ID != "" {
		staffSelectStatementParams = append(staffSelectStatementParams, ID)
		whereStatement = `where (` + column + `` + operator + `$1 and s."deleted_at" is null) and s."id"<>$2`
	}
	statement := `select distinct count(s."id") from "master_staffs" s ` + staffJoin + ` ` + whereStatement + ` ` + staffGroupBy
	err = repository.DB.QueryRow(statement, staffSelectStatementParams...).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}
