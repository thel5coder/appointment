package requests

type ClinicRequest struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	PICName     string `json:"pic_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}
