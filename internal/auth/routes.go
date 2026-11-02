package auth

import (
	"database/sql"
	"net/http"

	"github.com/Sabbir185/gpi/config"
	smsnet "github.com/Sabbir185/gpi/infra/message/sms_net"
	"github.com/Sabbir185/gpi/internal/device"
	"github.com/Sabbir185/gpi/internal/middleware"
	"github.com/Sabbir185/gpi/internal/user"
	"github.com/Sabbir185/gpi/pkg/token"
)

type Deps struct {
	DB           *sql.DB
	SMSClient    *smsnet.SMSClient
	Verifier     IDTokenVerifier
	Tokens       *token.Manager
	Queue        WelcomeQueue
	Authenticate middleware.Authenticator
	Cnf          *config.Config
}

func Routes(mux *http.ServeMux, deps *Deps) {
	// dependencies
	authService := NewAuthService(
		NewAuthRepository(deps.DB),
		user.NewUserRepository(deps.DB),
		device.NewDeviceRepository(deps.DB),
		deps.SMSClient,
		deps.Verifier,
		deps.Tokens,
		deps.Queue,
		deps.Cnf,
	)
	handler := NewAuthHandler(authService)
	// routes
	mux.HandleFunc("POST /auth/otp", handler.RequestOTP)
	mux.HandleFunc("POST /auth/otp/verify", handler.VerifyOTP)
	mux.HandleFunc("POST /auth/firebase", handler.SignInWithFirebase)
	mux.HandleFunc("POST /auth/refresh", handler.Refresh)
	mux.HandleFunc("POST /auth/sign-out", deps.Authenticate(handler.SignOut))
}
