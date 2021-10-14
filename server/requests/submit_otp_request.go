package requests

type SubmitOTPRequest struct {
	Type        string `json:"type"`
	MobilePhone string `json:"mobile_phone"`
	OTP         string `json:"otp"`
}
