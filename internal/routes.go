package internal

import (
	"net/http"

	"github.com/Sabbir185/gpi/internal/country"
	"github.com/Sabbir185/gpi/internal/healthz"
	"github.com/Sabbir185/gpi/internal/user"
	"github.com/Sabbir185/gpi/pkg/httpx"
)

func handleSideRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		httpx.SendSuccess(
			w,
			http.StatusOK,
			httpx.CodeDataFetch,
			"Welcome to the GPI Connect API 🚀",
		)
		return
	}
	// 404 Not Found
	httpx.SendError(
		w,
		http.StatusNotFound,
		httpx.CodeBadRequest,
		"The requested route was not found",
		nil,
	)
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
