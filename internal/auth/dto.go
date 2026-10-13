package auth

type SendRegistrationOTPRequest struct {
	Identity string `json:"identity" validate:"required,phone"`
	Type     string `json:"type" validate:"required,eq=register"`
}

type VerifyOTPRequest struct {
	Identity string `json:"identity" validate:"required,phone"`
	Type     string `json:"type" validate:"required,eq=register"`
	OTP      string `json:"otp" validate:"required,len=6"`
}
