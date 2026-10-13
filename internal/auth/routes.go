package auth

import (
	"database/sql"
	"net/http"
)

type Deps struct {
	DB *sql.DB
}

func Routes(mux *http.ServeMux, deps *Deps) {
	// dependencies
	authRepo := NewAuthRepository(deps.DB)
	authService := NewAuthService(authRepo)
	handler := NewAuthHandler(authService)
	// routes
	mux.HandleFunc("POST /auth/send-otp", handler.SendRegistrationOTP)
}
