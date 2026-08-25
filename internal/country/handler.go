package country

import (
	"net/http"

	"github.com/Sabbir185/gpc/pkg/response"
)

func AddNewCountry(w http.ResponseWriter, r *http.Request) {
	response.SendSuccess(w, http.StatusCreated, "Country added successfully", nil)
}
