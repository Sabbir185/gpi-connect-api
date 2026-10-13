package auth

import "net/http"

func Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/send-otp", SendRegistrationOTP)
}
