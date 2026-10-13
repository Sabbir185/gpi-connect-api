package auth

import (
	"time"

	"github.com/Sabbir185/gpi/entity"
	"github.com/Sabbir185/gpi/pkg/code"
)

type AuthService interface {
	SendOTP(payload *SendRegistrationOTPRequest) error
	VerifyOTP(payload *VerifyOTPRequest) error
}

type authService struct {
	authRepo AuthRepository
}

func NewAuthService(authRepo AuthRepository) AuthService {
	return &authService{authRepo: authRepo}
}

func (s *authService) SendOTP(payload *SendRegistrationOTPRequest) error {
	otp := code.GenerateRandomOTP(6)
	err := s.authRepo.CreateOTP(
		payload.Identity,
		otp,
		entity.OtpType(payload.Type),
		entity.OtpChannelSMS,
		1*time.Minute,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) VerifyOTP(payload *VerifyOTPRequest) error {
	return nil
}
