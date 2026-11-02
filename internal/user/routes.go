package user

import (
	"database/sql"
	"net/http"

	"github.com/Sabbir185/gpi/internal/designation"
	"github.com/Sabbir185/gpi/internal/device"
	"github.com/Sabbir185/gpi/internal/middleware"
)

type Deps struct {
	DB           *sql.DB
	Authenticate middleware.Authenticator
}

func Routes(mux *http.ServeMux, deps *Deps) {
	// dependencies
	userService := NewUserService(
		NewUserRepository(deps.DB),
		device.NewDeviceRepository(deps.DB),
		designation.NewDesignationRepository(deps.DB),
	)
	handler := NewUserHandler(userService)
	// routes
	mux.HandleFunc("GET /me", deps.Authenticate(handler.Me))
	mux.HandleFunc("PATCH /me", deps.Authenticate(handler.UpdateProfile))
	mux.HandleFunc("GET /me/devices", deps.Authenticate(handler.Devices))
	mux.HandleFunc("PUT /me/devices/{installationId}", deps.Authenticate(handler.RegisterDevice))
	mux.HandleFunc("DELETE /me/devices/{deviceId}", deps.Authenticate(handler.SignOutDevice))
}
