package designation

import (
	"database/sql"
	"net/http"

	"github.com/Sabbir185/gpi/internal/middleware"
)

type Deps struct {
	DB           *sql.DB
	Authenticate middleware.Authenticator
}

func Routes(mux *http.ServeMux, deps *Deps) {
	// dependencies
	handler := NewDesignationHandler(NewDesignationRepository(deps.DB))
	// routes
	mux.HandleFunc("GET /designations", deps.Authenticate(handler.Search))
}
