package auth

import (
	"context"
	"log"
	"time"

	"github.com/Sabbir185/gpi/entity"
	"github.com/Sabbir185/gpi/pkg/code"
)

type AuthService interface {
	SendOTP(ctx context.Context, payload *SendRegistrationOTPRequest) error
	VerifyOTP(payload *VerifyOTPRequest) error
}

type authService struct {
	authRepo AuthRepository
}

func NewAuthService(authRepo AuthRepository) AuthService {
	return &authService{authRepo: authRepo}
}

func (s *authService) SendOTP(ctx context.Context, payload *SendRegistrationOTPRequest) error {
	otp := code.GenerateRandomOTP(6)
	result, err := s.authRepo.CreateOTP(
		ctx,
		payload.Identity,
		otp,
		entity.OtpType(payload.Type),
		entity.OtpChannelSMS,
		1*time.Minute,
	)
	if err != nil {
		return err
	}
	log.Printf("OTP Created: %v", result)
	// TODO: send otp via whatsapp or sms
	return nil
}

func (s *authService) VerifyOTP(payload *VerifyOTPRequest) error {
	return nil
}
