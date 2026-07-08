package xerr

import "net/http"

const (
	CodeBadRequest   = 400001
	CodeUnauthorized = 401001
	CodeNotFound     = 404001
	CodeConflict     = 409001
	CodeInternal     = 500001
)

type CodeError struct {
	Status  int
	Code    int
	Message string
}

func (e *CodeError) Error() string {
	return e.Message
}

func NewBadRequest(message string) *CodeError {
	return &CodeError{
		Status:  http.StatusBadRequest,
		Code:    CodeBadRequest,
		Message: message,
	}
}

func NewUnauthorized(message string) *CodeError {
	return &CodeError{
		Status:  http.StatusUnauthorized,
		Code:    CodeUnauthorized,
		Message: message,
	}
}

func NewNotFound(message string) *CodeError {
	return &CodeError{
		Status:  http.StatusNotFound,
		Code:    CodeNotFound,
		Message: message,
	}
}

func NewConflict(message string) *CodeError {
	return &CodeError{
		Status:  http.StatusConflict,
		Code:    CodeConflict,
		Message: message,
	}
}

func NewInternal(message string) *CodeError {
	return &CodeError{
		Status:  http.StatusInternalServerError,
		Code:    CodeInternal,
		Message: message,
	}
}
