package viewmodel

type ClinicVm struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PICName     string `json:"pic_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
