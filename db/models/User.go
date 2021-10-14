package models

import "database/sql"

type User struct {
	ID                 string         `db:"id"`
	Name               string         `db:"name"`
	Email              string         `db:"email"`
	MobilePhone        string         `db:"mobile_phone"`
	Password           string         `db:"password"`
	IsActive           bool           `db:"is_active"`
	FcmDeviceToken     sql.NullString `db:"fcm_device_token"`
	ProfilePictureID   sql.NullString `db:"profile_picture_id"`
	ActivatedAt        sql.NullString `db:"activated_at"`
	CreatedAt          string         `db:"created_at"`
	UpdatedAt          string         `db:"updated_at"`
	DeletedAt          sql.NullString `db:"deleted_at"`
	Role               Role           `db:"role"`
	ProfilePicturePath sql.NullString `db:"profile_picture_path"`
}
