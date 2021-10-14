package requests

type ActivationCustomerRequest struct {
	OTP string `json:"otp" validate:"required"`
}
