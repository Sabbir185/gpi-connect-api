package auth

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Sabbir185/gpi/config"
	"github.com/Sabbir185/gpi/entity"
	smsnet "github.com/Sabbir185/gpi/infra/message/sms_net"
	"github.com/Sabbir185/gpi/pkg/code"
)

type AuthService interface {
	SendOTP(ctx context.Context, payload *SendRegistrationOTPRequest) (string, error)
	VerifyOTP(payload *VerifyOTPRequest) error
}

type authService struct {
	authRepo  AuthRepository
	smsClient *smsnet.SMSClient
	Cnf       *config.Config
}

func NewAuthService(authRepo AuthRepository, smsClient *smsnet.SMSClient, cnf *config.Config) AuthService {
	return &authService{
		authRepo:  authRepo,
		smsClient: smsClient,
		Cnf:       cnf,
	}
}

func (s *authService) SendOTP(ctx context.Context, payload *SendRegistrationOTPRequest) (string, error) {
	expireIn, err := strconv.Atoi(s.Cnf.SMSNet.ExpireInMin)
	if err != nil {
		expireIn = 2
	}
	otp := code.GenerateRandomOTP(6)

	result, err := s.authRepo.CreateOTP(
		ctx,
		payload.Identity,
		otp,
		entity.OtpType(payload.Type),
		entity.OtpChannelSMS,
		time.Minute * time.Duration(expireIn), // expire minutes dynamic
	)
	if err != nil {
		return "", err
	}
	log.Printf("OTP Created: %v", result)

	// if development, send otp in response
	if s.Cnf.App.Env == "development" {
		return result.OTP, nil
	}

	// real sms send for production
	smsResponse, err := s.smsClient.SendOTP(ctx, &smsnet.SendSMSPayload{
		To:  result.Identity,
		Msg: fmt.Sprintf("Your verification OTP is: %s. It will expire in %d minute.", otp, expireIn),
	})
	if err != nil {
		log.Printf("Failed to send OTP: %v", err)
		return "", err
	}
	log.Printf("OTP Sent: %v", smsResponse)
	return "", nil
}

func (s *authService) VerifyOTP(payload *VerifyOTPRequest) error {
	return nil
}
