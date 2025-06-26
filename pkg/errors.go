package pkg

import (
	"fmt"
	"net/http"
)

type ServiceError struct {
	Code          int
	ResponseCode  int
	Message       string
	InternalError error
}

func (e ServiceError) Error() string {
	if e.InternalError != nil {
		return fmt.Sprintf("error: {Message: %s, InternalError: %s}", e.Message, e.InternalError.Error())
	}

	return fmt.Sprintf("error: {Message: %s}", e.Message)

}

const (
	ErrBadRequest = 1 + iota
	ErrNotFound
	ErrInternalServer
)

var ServiceErrors = map[int]ServiceError{
	ErrBadRequest: {
		Code:         ErrBadRequest,
		ResponseCode: http.StatusBadRequest,
		Message:      "Bad request",
	},
	ErrNotFound: {
		Code:         ErrNotFound,
		ResponseCode: http.StatusNotFound,
		Message:      "Not found",
	},
	ErrInternalServer: {
		Code:         ErrInternalServer,
		ResponseCode: http.StatusInternalServerError,
		Message:      "Internal server error",
	},
}
