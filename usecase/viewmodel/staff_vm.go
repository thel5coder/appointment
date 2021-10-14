package viewmodel

type StaffVm struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	User        UserStaffVm     `json:"user"`
	Clinics     []StaffClinicVm `json:"clinics"`
}

type UserStaffVm struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	MobilePhone    string `json:"mobile_phone"`
	RoleID         string `json:"role_id"`
	ProfilePicture FileVm `json:"profile_picture"`
	IsActive       bool   `json:"is_active"`
}

type StaffClinicVm struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PICName     string `json:"pic_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}
