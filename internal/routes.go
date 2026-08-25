package internal

import (
	"net/http"

	"github.com/Sabbir185/gpc/internal/country"
	"github.com/Sabbir185/gpc/internal/healthz"
	"github.com/Sabbir185/gpc/internal/user"
	"github.com/Sabbir185/gpc/pkg/response"
)

func handleSideRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		response.SendSuccess(w, http.StatusOK, "Welcome to the GeoPunch Connect API 🚀", nil)
		return
	}
	// 404 Not Found
	response.SendError(w, http.StatusNotFound, "The requested route was not found", nil)
}

// RegisterRoutes will register all the routes
func RegisterRoutes(mux *http.ServeMux) *http.ServeMux {
	// ==================== API v1 ====================
	mux_v1 := http.NewServeMux()
	healthz.Routes(mux_v1)
	user.Routes(mux_v1)
	country.Routes(mux_v1)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", mux_v1))

	// ==================== Root / Fallback ====================
	mux.HandleFunc("/", handleSideRequest)

	return mux
}
