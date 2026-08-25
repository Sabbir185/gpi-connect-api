package response

import (
	"encoding/json"
	"net/http"
)

type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

type JsonResponse struct {
	Status  Status `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func SendSuccess(w http.ResponseWriter, code int, message string, data any) {
	response := &JsonResponse{
		Status:  StatusSuccess,
		Code:    code,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func SendError(w http.ResponseWriter, code int, message string, err any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&JsonResponse{
		Status:  StatusError,
		Code:    code,
		Message: message,
		Error:   err,
	})
}
