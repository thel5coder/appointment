package requests

type EditProfileRequest struct {
	Name             string `json:"name" validate:"required"`
	Email            string `json:"email" validate:"required,email"`
	MobilePhone      string `json:"mobile_phone" validate:"required"`
	BirthDate        string `json:"birth_date"`
	Password         string `json:"password"`
	ProfilePictureID string `json:"profile_picture_id"`
}
