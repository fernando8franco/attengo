package apperr

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mattn/go-sqlite3"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

func NewBadRequest(message string) *ErrorResponse {
	return &ErrorResponse{Code: http.StatusBadRequest, Message: message}
}

func NewUnauthorizedRequest(message string) *ErrorResponse {
	return &ErrorResponse{Code: http.StatusUnauthorized, Message: message}
}

type NotFoundError struct {
	Resource string
	ID       any
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id %v not found", e.Resource, e.ID)
}

type ConflictError struct {
	Resource string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("%s already exists", e.Resource)
}

type UnauthorizedError struct {
	Field   string
	Message string
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("validation failed on '%s': %s", e.Field, e.Message)
}

func IsUniqueConstraint(err error) bool {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
	}
	return false
}
