package models

type Clinic struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Address     string `db:"address"`
	PICName     string `db:"pic_name"`
	PhoneNumber string `db:"phone_number"`
	Email       string `db:"email"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}
