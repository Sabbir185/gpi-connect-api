package country

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Sabbir185/gpc/pkg/response"
	"github.com/Sabbir185/gpc/pkg/validation"
	"github.com/go-playground/validator/v10"
)

func AddNewCountry(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// payload from body with Unmarshal
	var body *CountryPayload
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validation check
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		formattedError := validation.FormatValidationErrors(err)
		response.SendError(w, http.StatusBadRequest, "Validation Error", formattedError)
		return
	}
	// service call
	response.SendSuccess(w, http.StatusCreated, "Country added successfully", nil)
}
