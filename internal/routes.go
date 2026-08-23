package internal

import (
	"net/http"

	"github.com/Sabbir185/gpc/internal/healthz"
)

// RegisterRoutes will register all the routes
func RegisterRoutes(mux *http.ServeMux) *http.ServeMux {
	healthz.Routes(mux)
	return mux
}
