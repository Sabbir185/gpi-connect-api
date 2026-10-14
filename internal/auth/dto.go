package auth

type SendRegistrationOTPRequest struct {
	Identity string `json:"identity" validate:"required,phone"`
	Type     string `json:"type" validate:"required,eq=register"`
}

type VerifyOTPRequest struct {
	Identity string `json:"identity" validate:"required,phone"`
	OTP      string `json:"otp" validate:"required,len=6"`
	Type     string `json:"type" validate:"required,eq=register"`
	FCMToken string `json:"fcm_token" validate:"required"`
	Device   struct {
		InstallationID string `json:"installation_id" validate:"required"`
		Platform       string `json:"platform" validate:"required"`
		OSVersion      string `json:"os_version" validate:"required"`
		AppVersion     string `json:"app_version" validate:"required"`
		Locale         string `json:"locale" validate:"required"`
		TimeZone       string `json:"timezone" validate:"required"`
	}
}
