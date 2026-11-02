package device

type DeviceRequest struct {
	InstallationID string  `json:"installationId" validate:"required,max=100"`
	Platform       string  `json:"platform" validate:"required,oneof=android ios web"`
	OSVersion      string  `json:"osVersion" validate:"max=50"`
	AppVersion     *string `json:"appVersion" validate:"omitempty,max=50"`
	Locale         string  `json:"locale" validate:"max=35"`
	Timezone       *string `json:"timezone" validate:"omitempty,max=64"`
	PushToken      *string `json:"pushToken" validate:"omitempty,max=4096"`
}

type DeviceResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Platform     string `json:"platform"`
	LastActiveAt string `json:"lastActiveAt"`
	Current      bool   `json:"current"`
}
