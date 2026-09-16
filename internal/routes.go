package internal

import (
	"database/sql"
	"net/http"

	"github.com/Sabbir185/gpi/config"
	"github.com/Sabbir185/gpi/internal/country"
	"github.com/Sabbir185/gpi/internal/healthz"
	"github.com/Sabbir185/gpi/internal/middleware"
	"github.com/Sabbir185/gpi/internal/user"
	"github.com/Sabbir185/gpi/pkg/httpx"
	"github.com/redis/go-redis/v9"
)

const apiV1Prefix = "/api/v1"

type Dependencies struct {
	Cnf   *config.Config
	DB    *sql.DB
	Redis *redis.Client
}

func RegisterRoutes(deps *Dependencies) http.Handler {
	rootMux := http.NewServeMux()
	v1Mux := http.NewServeMux()

	// Root & Health Check
	rootMux.HandleFunc("GET /{$}", rootHandler)
	rootMux.HandleFunc("GET /healthz", healthz.Healthz)

	// API V1 Routes
	registerV1Routes(v1Mux, deps)

	// Mount /api/v1
	rootMux.Handle(
		apiV1Prefix+"/",
		http.StripPrefix(apiV1Prefix, v1Mux),
	)

	// Global 404
	rootMux.HandleFunc("/", notFoundHandler)

	// Global Middlewares
	return middleware.RequestId(rootMux)
}

func registerV1Routes(mux *http.ServeMux, deps *Dependencies) {
	country.Routes(mux)
	user.Routes(mux)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	httpx.SendSuccess(
		w,
		http.StatusOK,
		httpx.CodeDataFetch,
		"Welcome to the GPI Connect API 🚀",
	)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	httpx.SendError(
		w,
		http.StatusNotFound,
		httpx.CodeNotFound,
		"The requested route was not found",
		nil,
	)
}
