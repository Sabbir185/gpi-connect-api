package healthz

import "net/http"

func Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", Healthz)
}
