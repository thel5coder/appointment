package requests

type LoginRequest struct {
	Email          string `json:"email"`
	MobilePhone    string `json:"mobile_phone"`
	Password       string `json:"password" validate:"required"`
	FcmDeviceToken string `json:"fcm_device_token"`
}
