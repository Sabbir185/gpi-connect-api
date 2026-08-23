package healthz

import (
	"encoding/json"
	"log"
	"net/http"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := map[string]string{
		"status": "okay",
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Error encoding response: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}
}
