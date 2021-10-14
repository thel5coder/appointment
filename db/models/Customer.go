package models

import (
	"database/sql"
	"time"
)

type Customer struct {
	ID                 string         `db:"id"`
	Name               string         `db:"name"`
	Sex                sql.NullString `db:"sex"`
	Address            sql.NullString `db:"address"`
	BirthDate          time.Time      `db:"birth_date"`
	MaritalStatus      sql.NullString `db:"marital_status"`
	PhoneNumber        sql.NullString `db:"phone_number"`
	MobilePhoneNumber1 sql.NullString `db:"mobile_phone_number_1"`
	MobilePhoneNumber2 sql.NullString `db:"mobile_phone_number_2"`
	Religion           sql.NullString `db:"religion"`
	Education          sql.NullString `db:"education"`
	Hobby              sql.NullString `db:"hobby"`
	Profession         sql.NullString `db:"profession"`
	Reference          sql.NullString `db:"reference"`
	Notes              sql.NullString `db:"notes"`
	CityID             sql.NullString `db:"city_id"`
	UserID             string         `db:"user_id"`
	CreatedAt          string         `db:"created_at"`
	UpdatedAt          string         `db:"updated_at"`
	DeletedAt          sql.NullString `db:"deleted_at"`
	User               User           `db:"user"`
}
