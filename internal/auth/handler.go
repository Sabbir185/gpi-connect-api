package auth

import (
	"net/http"

	"github.com/Sabbir185/gpi/internal/middleware"
	"github.com/Sabbir185/gpi/pkg/httpx"
)

type authHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *authHandler {
	return &authHandler{authService: authService}
}

func (h *authHandler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var payload RequestOTPRequest
	if !httpx.DecodeAndValidate(w, r, &payload) {
		return
	}
	result, err := h.authService.RequestOTP(r.Context(), &payload)
	if err != nil {
		httpx.SendAppError(w, r, err)
		return
	}
	httpx.SendSuccess(w, http.StatusOK, httpx.CodeOTPSent, "OTP sent successfully", result)
}

func (h *authHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var payload VerifyOTPRequest
	if !httpx.DecodeAndValidate(w, r, &payload) {
		return
	}
	result, err := h.authService.VerifyOTP(r.Context(), &payload)
	if err != nil {
		httpx.SendAppError(w, r, err)
		return
	}
	httpx.SendSuccess(w, http.StatusOK, httpx.CodeSignedIn, "Signed in successfully", result)
}

func (h *authHandler) SignInWithFirebase(w http.ResponseWriter, r *http.Request) {
	var payload FirebaseSignInRequest
	if !httpx.DecodeAndValidate(w, r, &payload) {
		return
	}
	result, err := h.authService.SignInWithFirebase(r.Context(), &payload)
	if err != nil {
		httpx.SendAppError(w, r, err)
		return
	}
	httpx.SendSuccess(w, http.StatusOK, httpx.CodeSignedIn, "Signed in successfully", result)
}

func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var payload RefreshRequest
	if !httpx.DecodeAndValidate(w, r, &payload) {
		return
	}
	result, err := h.authService.Refresh(r.Context(), payload.RefreshToken)
	if err != nil {
		httpx.SendAppError(w, r, err)
		return
	}
	httpx.SendSuccess(w, http.StatusOK, httpx.CodeTokenRefreshed, "Session refreshed", result)
}

func (h *authHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	auth := middleware.AuthFromContext(r.Context())
	if err := h.authService.SignOut(r.Context(), auth.UserID, auth.SessionID); err != nil {
		httpx.SendAppError(w, r, err)
		return
	}
	httpx.SendSuccess(w, http.StatusOK, httpx.CodeSignedOut, "Signed out successfully")
}
