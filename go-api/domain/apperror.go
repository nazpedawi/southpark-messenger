package domain

import "net/http"

// AppError is an error that also knows which HTTP status code fits it.
type AppError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

// AsMessage returns the human-readable part of the error.
func (e *AppError) AsMessage() string {
	return e.Message
}

// ValidationError: the client sent something wrong -> 400 Bad Request.
func ValidationError(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

// UnexpectedError: something broke on our side -> 500 Internal Server Error.
func UnexpectedError(message string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: message}
}
