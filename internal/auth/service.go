package auth

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	firebaseauth "firebase.google.com/go/v4/auth"
	"github.com/Sabbir185/gpi/config"
	"github.com/Sabbir185/gpi/entity"
	smsnet "github.com/Sabbir185/gpi/infra/message/sms_net"
	"github.com/Sabbir185/gpi/infra/queue"
	"github.com/Sabbir185/gpi/internal/device"
	"github.com/Sabbir185/gpi/internal/user"
	"github.com/Sabbir185/gpi/pkg/code"
	"github.com/Sabbir185/gpi/pkg/httpx"
	"github.com/Sabbir185/gpi/pkg/token"
	"github.com/Sabbir185/gpi/pkg/validation"
	"github.com/google/uuid"
)

const (
	otpLength         = 6
	otpMaxAttempts    = 5
	otpResendAfter    = 60 * time.Second
	otpHourlyLimit    = 5
	defaultOTPExpiry  = 3
	refreshReuseGrace = 60 * time.Second
	googleProvider    = "google.com"
)

type IDTokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*firebaseauth.Token, error)
}

type WelcomeQueue interface {
	EnqueueWelcomeEmail(payload queue.WelcomeEmailPayload) error
	EnqueueWelcomePush(payload queue.WelcomePushPayload) error
}

type AuthService interface {
	RequestOTP(ctx context.Context, payload *RequestOTPRequest) (*OTPRequestResponse, error)
	VerifyOTP(ctx context.Context, payload *VerifyOTPRequest) (*AuthResponse, error)
	SignInWithFirebase(ctx context.Context, payload *FirebaseSignInRequest) (*AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*SessionResponse, error)
	SignOut(ctx context.Context, userID int64, sessionID string) error
}

type authService struct {
	authRepo   AuthRepository
	userRepo   user.UserRepository
	deviceRepo device.DeviceRepository
	smsClient  *smsnet.SMSClient
	verifier   IDTokenVerifier
	tokens     *token.Manager
	queue      WelcomeQueue
	Cnf        *config.Config
}

func NewAuthService(
	authRepo AuthRepository,
	userRepo user.UserRepository,
	deviceRepo device.DeviceRepository,
	smsClient *smsnet.SMSClient,
	verifier IDTokenVerifier,
	tokens *token.Manager,
	queue WelcomeQueue,
	cnf *config.Config,
) AuthService {
	return &authService{
		authRepo:   authRepo,
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
		smsClient:  smsClient,
		verifier:   verifier,
		tokens:     tokens,
		queue:      queue,
		Cnf:        cnf,
	}
}

func (s *authService) RequestOTP(ctx context.Context, payload *RequestOTPRequest) (*OTPRequestResponse, error) {
	phone, err := validation.NormalizePhone(payload.Phone)
	if err != nil {
		return nil, errInvalidPhone()
	}

	latest, err := s.authRepo.FindLatestOTP(ctx, phone, entity.OtpTypeLogin)
	if err != nil {
		return nil, err
	}
	if latest != nil {
		if wait := otpResendAfter - time.Since(latest.CreatedAt); wait > 0 {
			return nil, errRateLimited(wait)
		}
	}

	sentLastHour, err := s.authRepo.CountOTPsSince(ctx, phone, entity.OtpTypeLogin, time.Now().Add(-time.Hour))
	if err != nil {
		return nil, err
	}
	if sentLastHour >= otpHourlyLimit {
		return nil, errRateLimited(time.Hour)
	}

	otpCode, err := code.GenerateRandomNumber(otpLength)
	if err != nil {
		return nil, err
	}
	expireIn, err := strconv.Atoi(s.Cnf.SMSNet.ExpireInMin)
	if err != nil || expireIn <= 0 {
		expireIn = defaultOTPExpiry
	}

	// WhatsApp delivery isn't wired yet, so every code goes out by SMS and the app shows that channel.
	otp, err := s.authRepo.CreateOTP(
		ctx,
		phone,
		code.HashOTP(s.Cnf.Auth.OTPSecret, phone, otpCode),
		entity.OtpTypeLogin,
		entity.OtpChannelSMS,
		time.Minute*time.Duration(expireIn),
	)
	if err != nil {
		return nil, err
	}

	if s.Cnf.App.Env == "development" {
		log.Printf("[dev] OTP for %s: %s", phone, otpCode)
	}

	if s.Cnf.SMSNet.Enabled {
		_, err = s.smsClient.SendOTP(ctx, &smsnet.SendSMSPayload{
			To:  phone,
			Msg: fmt.Sprintf("Your %s verification code is %s. It expires in %d minutes. Never share this code.", s.Cnf.App.Name, otpCode, expireIn),
		})
		if err != nil {
			_ = s.authRepo.ExpireOTP(context.WithoutCancel(ctx), otp.ID)
			return nil, fmt.Errorf("send otp sms: %w", err)
		}
	}

	return &OTPRequestResponse{
		RequestID:   otp.RequestID,
		Channel:     string(entity.OtpChannelSMS),
		ResendAfter: int(otpResendAfter.Seconds()),
		CodeLength:  otpLength,
	}, nil
}

func (s *authService) VerifyOTP(ctx context.Context, payload *VerifyOTPRequest) (*AuthResponse, error) {
	phone, err := validation.NormalizePhone(payload.Phone)
	if err != nil {
		return nil, errInvalidPhone()
	}
	otpID, codeHash, err := s.authRepo.ConsumeOTPAttempt(ctx, payload.RequestID, phone, entity.OtpTypeLogin, otpMaxAttempts)
	if err != nil {
		return nil, err
	}
	if codeHash == "" {
		return nil, errOTPExpired()
	}
	if !code.EqualHash(codeHash, code.HashOTP(s.Cnf.Auth.OTPSecret, phone, payload.Code)) {
		return nil, httpx.NewAppError(http.StatusUnprocessableEntity, httpx.CodeInvalidOTP, "The code is incorrect", nil)
	}

	consumed, err := s.authRepo.MarkOTPConsumed(ctx, otpID)
	if err != nil {
		return nil, err
	}
	if !consumed {
		return nil, errOTPExpired()
	}

	account, created, err := s.findOrCreatePhoneUser(ctx, phone)
	if err != nil {
		return nil, err
	}

	result, err := s.startSession(ctx, account, &payload.Device)
	if err != nil {
		return nil, err
	}
	if created {
		s.enqueueWelcome(account, payload.Device.Locale)
	}
	return result, nil
}

func (s *authService) findOrCreatePhoneUser(ctx context.Context, phone string) (*entity.User, bool, error) {
	account, err := s.userRepo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, false, err
	}
	if account != nil {
		if account.PhoneVerifiedAt == nil {
			if err := s.userRepo.MarkPhoneVerified(ctx, account.ID); err != nil {
				return nil, false, err
			}
		}
		return account, false, nil
	}

	now := time.Now()
	account, err = s.userRepo.Create(ctx, &user.NewUser{Phone: &phone, PhoneVerifiedAt: &now})
	if err != nil {
		return nil, false, err
	}
	if account == nil {
		account, err = s.userRepo.FindByPhone(ctx, phone)
		return account, false, err
	}
	return account, true, nil
}

func (s *authService) SignInWithFirebase(ctx context.Context, payload *FirebaseSignInRequest) (*AuthResponse, error) {
	idToken, err := s.verifier.VerifyIDToken(ctx, payload.IDToken)
	if err != nil {
		log.Printf("firebase id token rejected: %v", err)
		return nil, errGoogleFailed()
	}
	if idToken.Firebase.SignInProvider != googleProvider {
		return nil, errGoogleFailed()
	}

	identity := firebaseIdentity(idToken, &payload.Profile)
	account, created, err := s.findOrCreateFirebaseUser(ctx, identity)
	if err != nil {
		return nil, err
	}

	result, err := s.startSession(ctx, account, &payload.Device)
	if err != nil {
		return nil, err
	}
	if created {
		s.enqueueWelcome(account, payload.Device.Locale)
	}
	return result, nil
}

type verifiedIdentity struct {
	UID           string
	Email         *string
	EmailVerified bool
	Name          string
	Picture       *string
}

// firebaseIdentity trusts the verified token first; the app's profile only fills gaps for display fields.
func firebaseIdentity(idToken *firebaseauth.Token, profile *SocialProfile) *verifiedIdentity {
	identity := &verifiedIdentity{UID: idToken.UID}
	if email, _ := idToken.Claims["email"].(string); email != "" {
		email = strings.ToLower(email)
		identity.Email = &email
		identity.EmailVerified, _ = idToken.Claims["email_verified"].(bool)
	}
	identity.Name, _ = idToken.Claims["name"].(string)
	if identity.Name == "" && profile.FullName != nil {
		identity.Name = strings.TrimSpace(*profile.FullName)
	}
	if picture, _ := idToken.Claims["picture"].(string); picture != "" {
		identity.Picture = &picture
	}
	if profile.AvatarURL != nil {
		identity.Picture = profile.AvatarURL
	}
	return identity
}

func (s *authService) findOrCreateFirebaseUser(ctx context.Context, identity *verifiedIdentity) (*entity.User, bool, error) {
	account, err := s.userRepo.FindByFirebaseUID(ctx, identity.UID)
	if err != nil || account != nil {
		return account, false, err
	}

	// Only a verified email may claim an existing account, otherwise anyone could take it over.
	if identity.Email != nil && identity.EmailVerified {
		account, err = s.userRepo.FindByEmail(ctx, *identity.Email)
		if err != nil {
			return nil, false, err
		}
		if account != nil {
			if err := s.userRepo.LinkFirebase(ctx, account.ID, identity.UID, true); err != nil {
				return nil, false, err
			}
			return account, false, nil
		}
	}

	newUser := &user.NewUser{
		Name:        identity.Name,
		ProfileImg:  identity.Picture,
		FirebaseUID: &identity.UID,
	}
	if identity.Email != nil && identity.EmailVerified {
		now := time.Now()
		newUser.Email = identity.Email
		newUser.EmailVerifiedAt = &now
	}

	account, err = s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, false, err
	}
	if account == nil {
		account, err = s.userRepo.FindByFirebaseUID(ctx, identity.UID)
		if err == nil && account == nil {
			return nil, false, httpx.NewAppError(http.StatusConflict, httpx.CodeConflict, "This email is already linked to another account", nil)
		}
		return account, false, err
	}
	return account, true, nil
}

func (s *authService) startSession(ctx context.Context, account *entity.User, deviceInfo *device.DeviceRequest) (*AuthResponse, error) {
	if account.Status != entity.UserStatusActive {
		return nil, httpx.NewAppError(http.StatusForbidden, httpx.CodeForbidden, "This account is not active", nil)
	}

	deviceID, err := s.deviceRepo.Upsert(ctx, account.ID, deviceInfo, true)
	if err != nil {
		return nil, err
	}
	// One live session per install: signing in again replaces the old one.
	if _, err := s.deviceRepo.RevokeSessions(ctx, account.ID, deviceID); err != nil {
		return nil, err
	}

	refreshToken, err := token.NewOpaqueToken()
	if err != nil {
		return nil, err
	}
	session := &entity.Session{
		ID:               uuid.NewString(),
		UserID:           account.ID,
		DeviceID:         &deviceID,
		RefreshTokenHash: token.Hash(refreshToken),
		ExpiresAt:        time.Now().Add(s.Cnf.Auth.RefreshTokenTTL),
	}
	if err := s.authRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	if err := s.userRepo.TouchLastLogin(ctx, account.ID); err != nil {
		log.Printf("touch last login for user %d: %v", account.ID, err)
	}

	sessionResponse, err := s.issueSessionResponse(account.ID, session.ID, refreshToken)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{Session: sessionResponse, User: user.ToUserResponse(account)}, nil
}

func (s *authService) issueSessionResponse(userID int64, sessionID, refreshToken string) (*SessionResponse, error) {
	accessToken, expiresAt, err := s.tokens.IssueAccessToken(userID, sessionID)
	if err != nil {
		return nil, err
	}
	return &SessionResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.UnixMilli(),
	}, nil
}

// Refresh rotates the refresh token on every use. Presenting a token that was
// already rotated away means it leaked, so the whole session is revoked, except
// within a short grace window that covers a retried request whose response was lost.
func (s *authService) Refresh(ctx context.Context, refreshToken string) (*SessionResponse, error) {
	presentedHash := token.Hash(refreshToken)

	session, err := s.authRepo.FindSessionByRefreshHash(ctx, presentedHash)
	if err != nil {
		return nil, err
	}
	if session == nil {
		session, err = s.authRepo.FindSessionByPreviousRefreshHash(ctx, presentedHash)
		if err != nil {
			return nil, err
		}
		if session == nil {
			return nil, errUnauthorized()
		}
		if session.RotatedAt == nil || time.Since(*session.RotatedAt) > refreshReuseGrace {
			log.Printf("refresh token reuse detected, revoking session %s", session.ID)
			if err := s.authRepo.RevokeSession(ctx, session.ID); err != nil {
				return nil, err
			}
			return nil, errUnauthorized()
		}
	}

	if session.RevokedAt != nil || time.Now().After(session.ExpiresAt) {
		return nil, errUnauthorized()
	}
	account, err := s.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}
	if account == nil || account.Status != entity.UserStatusActive {
		return nil, errUnauthorized()
	}

	nextRefreshToken, err := token.NewOpaqueToken()
	if err != nil {
		return nil, err
	}
	rotated, err := s.authRepo.RotateSession(
		ctx,
		session.ID,
		presentedHash,
		token.Hash(nextRefreshToken),
		time.Now().Add(s.Cnf.Auth.RefreshTokenTTL),
	)
	if err != nil {
		return nil, err
	}
	if !rotated {
		return nil, errUnauthorized()
	}

	return s.issueSessionResponse(session.UserID, session.ID, nextRefreshToken)
}

func (s *authService) SignOut(ctx context.Context, userID int64, sessionID string) error {
	session, err := s.authRepo.FindSession(ctx, sessionID, userID)
	if err != nil || session == nil {
		return err
	}
	if err := s.authRepo.RevokeSession(ctx, session.ID); err != nil {
		return err
	}
	if session.DeviceID != nil {
		return s.deviceRepo.ClearPushToken(ctx, *session.DeviceID)
	}
	return nil
}

// enqueueWelcome never fails the sign-in: the account already exists and the user is waiting.
func (s *authService) enqueueWelcome(account *entity.User, locale string) {
	if account.Email != nil {
		err := s.queue.EnqueueWelcomeEmail(queue.WelcomeEmailPayload{
			UserID: account.ID,
			Email:  *account.Email,
			Name:   account.Name,
			Locale: locale,
		})
		if err != nil {
			log.Printf("enqueue welcome email for user %d: %v", account.ID, err)
		}
	}
	if err := s.queue.EnqueueWelcomePush(queue.WelcomePushPayload{UserID: account.ID}); err != nil {
		log.Printf("enqueue welcome push for user %d: %v", account.ID, err)
	}
}

func errInvalidPhone() error {
	return httpx.NewAppError(http.StatusBadRequest, httpx.CodeInvalidPhone, "Invalid phone number", nil)
}

func errOTPExpired() error {
	return httpx.NewAppError(http.StatusGone, httpx.CodeOTPExpired, "The code has expired, request a new one", nil)
}

func errGoogleFailed() error {
	return httpx.NewAppError(http.StatusUnauthorized, httpx.CodeGoogleFailed, "Google sign-in could not be verified", nil)
}

func errUnauthorized() error {
	return httpx.NewAppError(http.StatusUnauthorized, httpx.CodeUnauthorized, "Session expired, sign in again", nil)
}

func errRateLimited(wait time.Duration) error {
	return httpx.NewAppError(
		http.StatusTooManyRequests,
		httpx.CodeRateLimited,
		"Too many attempts, try again later",
		map[string]int{"retryAfter": int(math.Ceil(wait.Seconds()))},
	)
}
