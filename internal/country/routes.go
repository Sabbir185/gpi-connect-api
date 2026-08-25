package country

import "net/http"

func Routes(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("POST /countries", AddNewCountry)
	return mux
}
