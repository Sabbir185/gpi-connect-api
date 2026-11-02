package httpx

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Sabbir185/gpi/pkg/validation"
)

const maxBodyBytes = 1 << 20

// DecodeAndValidate writes the error response itself and reports whether the handler may continue.
func DecodeAndValidate(w http.ResponseWriter, r *http.Request, payload any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		SendError(w, http.StatusBadRequest, CodeInvalidRequest, "Invalid request body", nil)
		return false
	}
	if err := validation.Struct(payload); err != nil {
		SendError(w, http.StatusBadRequest, CodeValidationError, "Validation error", validation.FormatValidationErrors(err))
		return false
	}
	return true
}

// SendAppError never leaks internal error text to clients; unexpected errors are logged instead.
func SendAppError(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		SendError(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Data)
		return
	}
	log.Printf("%s %s: %v", r.Method, r.URL.Path, err)
	SendError(w, http.StatusInternalServerError, CodeInternalError, "Internal server error", nil)
}
