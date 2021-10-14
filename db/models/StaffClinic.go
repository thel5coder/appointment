package models

type ClinicStaff struct {
	ID        string `db:"id"`
	StaffID   string `db:"staff_id"`
	ClinicID  string `db:"clinic_id"`
	CreatedAt string `db:"created_at"`
}
