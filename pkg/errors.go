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
		return fmt.Sprintf("error: %s", e.InternalError.Error())
	}

	return fmt.Sprintf("error: %s", e.Message)

}

const (
	ErrBadRequest = 1 + iota
	ErrNotFound
	ErrInternalServer
	ErrConflict
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
	ErrConflict: {
		Code:         ErrConflict,
		ResponseCode: http.StatusConflict,
		Message:      "Already exists",
	},
}
