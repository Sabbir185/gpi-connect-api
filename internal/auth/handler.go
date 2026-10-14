package auth

import (
	"encoding/json"
	"errors"
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
		httpx.SendError(w, http.StatusBadRequest, httpx.CodeInvalidRequest, "Invalid request", err.Error())
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
		var appErr *httpx.AppError
		if errors.As(err, &appErr) {
			// generic error
			var data any
			if h.Cnf.App.Env == "development" {
				data = appErr.Data
			}
			httpx.SendError(w, appErr.StatusCode, appErr.Code, appErr.Message, data)
			return
		}
		// fallback
		httpx.SendError(w, http.StatusInternalServerError, httpx.CodeInternalError, "Internal server error", err.Error())
		return
	}

	// if development, send otp in response
	var data any
	if h.Cnf.App.Env == "development" {
		data = map[string]string{
			"otp": otp,
			"env": h.Cnf.App.Env,
		}
	}

	httpx.SendSuccess(
		w,
		http.StatusOK,
		httpx.CodeOTPSent,
		"OTP sent successfully",
		data,
	)
}

func (h *authHandler) VerifyRegistrationOTP(w http.ResponseWriter, r *http.Request) {
	var payload *VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpx.SendError(w, http.StatusBadRequest, httpx.CodeInvalidRequest, "Invalid request", err.Error())
		return
	}

	validate := validator.New()
	_ = validate.RegisterValidation("phone", validation.ValidateGlobalPhone)
	if err := validate.Struct(payload); err != nil {
		formattedError := validation.FormatValidationErrors(err)
		httpx.SendError(w, http.StatusBadRequest, httpx.CodeValidationError, "Validation error", formattedError)
		return
	}

	// Handover to AuthService
	result, err := h.authService.VerifyRegistrationOTP(r.Context(), payload)
	if err != nil {
		var appErr *httpx.AppError
		if errors.As(err, &appErr) {
			httpx.SendError(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Data)
			return
		}
		// fallback
		httpx.SendError(w, http.StatusInternalServerError, httpx.CodeInternalError, "Internal server error", err.Error())
		return
	}

	httpx.SendSuccess(
		w,
		http.StatusOK,
		httpx.CodeVerified,
		"User verified successfully",
		result,
	)
}
