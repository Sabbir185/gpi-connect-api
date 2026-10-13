package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Sabbir185/gpi/pkg/httpx"
	"github.com/Sabbir185/gpi/pkg/validation"
	"github.com/go-playground/validator/v10"
)

type authHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *authHandler {
	return &authHandler{authService: authService}
}

func (h *authHandler) SendRegistrationOTP(w http.ResponseWriter, r *http.Request) {
	var payload SendRegistrationOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpx.SendError(w, http.StatusBadRequest, httpx.CodeInvalidRequest, "Invalid request", err)
		return
	}
	validate := validator.New()
	_ = validate.RegisterValidation("phone", validation.ValidateGlobalPhone)
	if err := validate.Struct(&payload); err != nil {
		formattedError := validation.FormatValidationErrors(err)
		httpx.SendError(w, http.StatusBadRequest, httpx.CodeValidationError, "Validation error", formattedError)
		return
	}
	err := h.authService.SendOTP(r.Context(), &payload)
	if err != nil {
		httpx.SendError(w, http.StatusInternalServerError, httpx.CodeInternalError, "Internal server error", err)
		return
	}
	httpx.SendSuccess(w, http.StatusOK, httpx.CodeOTPSent, "OTP sent successfully")
}
