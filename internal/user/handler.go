package user

import (
	"net/http"

	"github.com/Sabbir185/gpi/internal/device"
	"github.com/Sabbir185/gpi/internal/middleware"
	"github.com/Sabbir185/gpi/pkg/httpx"
)

type userHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *userHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) Me(w http.ResponseWriter, r *http.Request) {
	auth := middleware.AuthFromContext(r.Context())
	result, err := h.userService.Me(r.Context(), auth.UserID)
	if err != nil {
		httpx.SendAppError(w, r, err)
		return
	}
	httpx.SendSuccess(w, http.StatusOK, httpx.CodeDataFetch, "Profile fetched", result)
}

func (h *userHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var payload UpdateProfileRequest
	if !httpx.DecodeAndValidate(w, r, &payload) {
		return
	}
	auth := middleware.AuthFromContext(r.Context())
	result, err := h.userService.UpdateProfile(r.Context(), auth.UserID, &payload)
	if err != nil {
		httpx.SendAppError(w, r, err)
		return
	}
	httpx.SendSuccess(w, http.StatusOK, httpx.CodeDataUpdated, "Profile updated", result)
}

func (h *userHandler) Devices(w http.ResponseWriter, r *http.Request) {
	auth := middleware.AuthFromContext(r.Context())
	result, err := h.userService.Devices(r.Context(), auth.UserID, auth.SessionID)
	if err != nil {
		httpx.SendAppError(w, r, err)
		return
	}
	httpx.SendSuccess(w, http.StatusOK, httpx.CodeDataFetch, "Devices fetched", result)
}

func (h *userHandler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	var payload device.DeviceRequest
	if !httpx.DecodeAndValidate(w, r, &payload) {
		return
	}
	if payload.InstallationID != r.PathValue("installationId") {
		httpx.SendError(w, http.StatusBadRequest, httpx.CodeValidationError, "installationId does not match the path", nil)
		return
	}
	auth := middleware.AuthFromContext(r.Context())
	if err := h.userService.RegisterDevice(r.Context(), auth.UserID, &payload); err != nil {
		httpx.SendAppError(w, r, err)
		return
	}
	httpx.SendSuccess(w, http.StatusOK, httpx.CodeDataUpdated, "Device registered")
}

func (h *userHandler) SignOutDevice(w http.ResponseWriter, r *http.Request) {
	auth := middleware.AuthFromContext(r.Context())
	if err := h.userService.SignOutDevice(r.Context(), auth.UserID, r.PathValue("deviceId")); err != nil {
		httpx.SendAppError(w, r, err)
		return
	}
	httpx.SendSuccess(w, http.StatusOK, httpx.CodeSignedOut, "Device signed out")
}
