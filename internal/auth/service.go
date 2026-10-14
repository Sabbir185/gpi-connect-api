package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Sabbir185/gpi/config"
	"github.com/Sabbir185/gpi/entity"
	smsnet "github.com/Sabbir185/gpi/infra/message/sms_net"
	"github.com/Sabbir185/gpi/pkg/code"
	"github.com/Sabbir185/gpi/pkg/httpx"
)

type AuthService interface {
	SendOTP(ctx context.Context, payload *SendRegistrationOTPRequest) (string, error)
	VerifyRegistrationOTP(ctx context.Context, payload *VerifyOTPRequest) (any, error)
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
	// check otp already sent or not
	otp, err := s.authRepo.FindOTPByIdentityAndType(ctx, payload.Identity, entity.OtpType(payload.Type))
	if err != nil {
		fmt.Printf("%v", err)
	}
	if otp != nil {
		return "", httpx.NewAppError(
			http.StatusTooManyRequests,
			httpx.CodeTooManyRequests,
			"OTP already sent. Please try again later.",
			map[string]string{
				"otp": otp.OTP,
				"env": s.Cnf.App.Env,
			},
		)
	}

	// Generate new OTP and sending
	expireIn, err := strconv.Atoi(s.Cnf.SMSNet.ExpireInMin)
	if err != nil {
		expireIn = 2
	}
	otpStr := code.GenerateRandomNumber(6)

	result, err := s.authRepo.CreateOTP(
		ctx,
		payload.Identity,
		otpStr,
		entity.OtpType(payload.Type),
		entity.OtpChannelSMS,
		time.Minute*time.Duration(expireIn), // expire minutes dynamic
	)
	if err != nil {
		return "", err
	}

	// if development, send otp in response
	if s.Cnf.App.Env == "development" {
		return result.OTP, nil
	}

	// real sms send for production
	_, err = s.smsClient.SendOTP(ctx, &smsnet.SendSMSPayload{
		To:  result.Identity,
		Msg: fmt.Sprintf("Your verification OTP is: %s. It will expire in %d minute.", otpStr, expireIn),
	})
	if err != nil {
		log.Printf("Failed to send OTP: %v", err)
		return "", err
	}
	return "", nil
}

func (s *authService) VerifyRegistrationOTP(ctx context.Context, payload *VerifyOTPRequest) (any, error) {
	// check otp already expired or exits
	// verify otp
	// if user already exists, return error
	// if user doesn't exist, create user
	// if otp is valid, create user -- phone, fcm_token, verified_at, status etc
	// then return user and access token and refresh token
	return "", nil
}
