package requests

type CustomerRequest struct {
	Name             string `json:"name" validate:"required"`
	Sex              string `json:"sex" validate:"required"`
	BirthDate        string `json:"birth_date" validate:"required"`
	Email            string `json:"email" validate:"required"`
	MobilePhone      string `json:"mobile_phone" validate:"required"`
	Address          string `json:"address"`
	Password         string `json:"password"`
	ProfilePictureID string `json:"profile_picture_id"`
	IsActive         bool   `json:"is_active"`
}
