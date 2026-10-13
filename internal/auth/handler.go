package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Sabbir185/gpi/config"
	"github.com/Sabbir185/gpi/pkg/httpx"
	"github.com/Sabbir185/gpi/pkg/validation"
	"github.com/go-playground/validator/v10"
)

type authHandler struct {
	authService AuthService
	Cnf         *config.Config
}

func NewAuthHandler(authService AuthService, cnf *config.Config) *authHandler {
	return &authHandler{authService: authService, Cnf: cnf}
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
	otp, err := h.authService.SendOTP(r.Context(), &payload)
	if err != nil {
		httpx.SendError(w, http.StatusInternalServerError, httpx.CodeInternalError, "Internal server error", err)
		return
	}

	// if development, send otp in response
	var data = make(map[string]string)
	if h.Cnf.App.Env == "development" {
		data["otp"] = otp
		data["env"] = h.Cnf.App.Env
	}
	
	httpx.SendSuccess(
		w,
		http.StatusOK,
		httpx.CodeOTPSent,
		"OTP sent successfully",
		data,
	)
}
