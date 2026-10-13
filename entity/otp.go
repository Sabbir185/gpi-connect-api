package entity

import "time"

type Otp struct {
	ID        int64      `db:"id" json:"id"`
	Identity  string     `db:"identity" json:"identity"`
	OTP       string     `db:"otp" json:"otp"`
	Type      OtpType    `db:"type" json:"type"`
	Channel   OtpChannel `db:"channel" json:"channel"`
	Attempt   int        `db:"attempt" json:"attempt"`
	ExpiresAt time.Time  `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

type OtpType string

const (
	OtpTypeLogin    OtpType = "login"
	OtpTypeRegister OtpType = "register"
	OtpTypeReset    OtpType = "reset"
)

type OtpChannel string

const (
	OtpChannelSMS      OtpChannel = "sms"
	OtpChannelEmail    OtpChannel = "email"
	OtpChannelWhatsapp OtpChannel = "whatsapp"
	OtpChannelPush     OtpChannel = "push"
)
