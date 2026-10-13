package auth

import (
	"database/sql"
	"net/http"

	"github.com/Sabbir185/gpi/config"
	smsnet "github.com/Sabbir185/gpi/infra/message/sms_net"
)

type Deps struct {
	DB        *sql.DB
	SMSClient *smsnet.SMSClient
	Cnf       *config.Config
}

func Routes(mux *http.ServeMux, deps *Deps) {
	// dependencies
	authRepo := NewAuthRepository(deps.DB)
	authService := NewAuthService(authRepo, deps.SMSClient, deps.Cnf)
	handler := NewAuthHandler(authService, deps.Cnf)
	// routes
	mux.HandleFunc("POST /auth/send-otp", handler.SendRegistrationOTP)
}
