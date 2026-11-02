package auth

import (
	"github.com/Sabbir185/gpi/internal/device"
	"github.com/Sabbir185/gpi/internal/user"
)

type RequestOTPRequest struct {
	Phone   string `json:"phone" validate:"required,phone"`
	Channel string `json:"channel" validate:"required,oneof=sms whatsapp"`
}

type OTPRequestResponse struct {
	RequestID   string `json:"requestId"`
	Channel     string `json:"channel"`
	ResendAfter int    `json:"resendAfter"`
	CodeLength  int    `json:"codeLength"`
}

type VerifyOTPRequest struct {
	RequestID string               `json:"requestId" validate:"required,uuid"`
	Phone     string               `json:"phone" validate:"required,phone"`
	Code      string               `json:"code" validate:"required,len=6,numeric"`
	Device    device.DeviceRequest `json:"device" validate:"required"`
}

type SocialProfile struct {
	FullName  *string `json:"fullName" validate:"omitempty,max=100"`
	Email     *string `json:"email" validate:"omitempty,email"`
	AvatarURL *string `json:"avatarUrl" validate:"omitempty,url"`
}

type FirebaseSignInRequest struct {
	IDToken  string               `json:"idToken" validate:"required"`
	Provider string               `json:"provider" validate:"required,oneof=google"`
	Profile  SocialProfile        `json:"profile"`
	Device   device.DeviceRequest `json:"device" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type SessionResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int64  `json:"expiresAt"`
}

type AuthResponse struct {
	Session *SessionResponse   `json:"session"`
	User    *user.UserResponse `json:"user"`
}
