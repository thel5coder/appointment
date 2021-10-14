package viewmodel

type UserVm struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	MobilePhone    string `json:"mobile_phone"`
	FcmDeviceToken string `json:"fcm_device_token"`
	ProfilePicture FileVm `json:"profile_picture"`
	IsActive       bool   `json:"is_active"`
	ActivatedAt    string `json:"activated_at"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	Role           UserRoleVm `json:"role"`
}

type UserRoleVm struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}
