package apperr

import (
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

func NewForbiddenRequest(message string) *ErrorResponse {
	return &ErrorResponse{Code: http.StatusForbidden, Message: message}
}

func NewBadRequest(message string) *ErrorResponse {
	return &ErrorResponse{Code: http.StatusBadRequest, Message: message}
}

func NewUnauthorizedRequest(message string) *ErrorResponse {
	return &ErrorResponse{Code: http.StatusUnauthorized, Message: message}
}
