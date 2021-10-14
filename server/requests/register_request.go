package requests

type RegisterRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email"`
	MobilePhone string `json:"mobile_phone" validate:"required"`
	BirthDate   string `json:"birth_date" validate:"required"`
	Password    string `json:"password" validate:"required"`
}
