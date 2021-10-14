package requests

type OTPRequest struct {
	MobilePhone string `json:"mobile_phone" validate:"required"`
}
