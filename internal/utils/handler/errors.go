package handler

import "fmt"

type ErrBadRequest struct {
	message string
}

func (e ErrBadRequest) Error() string {
	return e.message
}

func NewErrBadRequest(message string) ErrBadRequest {
	return ErrBadRequest{
		message: message,
	}
}

type ErrForbidden struct {
	message string
}

func (e ErrForbidden) Error() string {
	return fmt.Sprintf("forbidden: %s", e.message)
}

func NewErrForbidden(message string) error {
	return ErrForbidden{message: message}
}

func mandatoryNotPresent(param, paramType string) error {
	return ErrBadRequest{
		message: fmt.Sprintf("%s %s is missing", paramType, param),
	}
}

func mandatoryQueryNotPresent(param string) error {
	return mandatoryNotPresent(param, "URL query parameter")
}
