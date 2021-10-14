package requests

type UserRequest struct {
	Name             string `json:"name" validate:"required"`
	Email            string `json:"email" validate:"required"`
	MobilePhone      string `json:"mobile_phone" validate:"required"`
	Password         string `json:"password"`
	RoleID           string `json:"role_id"`
	ProfilePictureID string `json:"profile_picture_id"`
	IsActive         bool   `json:"is_active"`
	ActivatedAt      string `json:"activated_at"`
}
