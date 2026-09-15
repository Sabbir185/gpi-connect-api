package httpx

import (
	"encoding/json"
	"net/http"
)

type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

type Code string

const (
	CodeDataFetch       Code = "data_fetched"
	CodeDataCreated     Code = "data_created"
	CodeDataUpdated     Code = "data_updated"
	CodeDataDeleted     Code = "data_deleted"
	CodeDataNotFound    Code = "data_not_found"
	CodeDataExists      Code = "data_already_exists"
	CodeInvalidRequest  Code = "invalid_request"
	CodeValidationError Code = "validation_error"
	CodeInternalError   Code = "internal_server_error"
	CodeNotFound        Code = "not_found"
	CodeBadRequest      Code = "bad_request"
	CodeForbidden       Code = "forbidden"
	CodeUnauthorized    Code = "unauthorized"
	CodeConflict        Code = "conflict"
	CodeTooManyRequests Code = "too_many_requests"
)

type Meta struct {
	Page        int64 `json:"page"`
	Limit       int64 `json:"limit"`
	Total       int64 `json:"total"`
	TotalPages  int64 `json:"total_pages"`
	HasNext     int64 `json:"has_next"`
	HasPrevious int64 `json:"has_previous"`
}

type JsonResponse struct {
	Status  Status `json:"status"`
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    Meta   `json:"meta,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func SendSuccess(w http.ResponseWriter, statusCode int, code Code, message string, values ...any) {
	var data any
	var meta Meta
	if len(values) > 0 {
		data = values[0]
	}
	if len(values) > 1 {
		if m, ok := values[1].(Meta); ok {
			meta = m
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&JsonResponse{
		Status:  StatusSuccess,
		Code:    code,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func SendError(w http.ResponseWriter, statusCode int, code Code, message string, err any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&JsonResponse{
		Status:  StatusError,
		Code:    code,
		Message: message,
		Error:   err,
	})
}
