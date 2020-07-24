package common

import (
	"strings"
)

// Error types
const (
	NoRowsAffected = 101

	UserNotAllowed    = 301
	IncorrectPassword = 302

	CityNotFound    = 401
	CommentNotFount = 402
	UserNotFound    = 403
	AirportNotFound = 404
	RouteNotFound   = 405
)

// GeneralError - General description error with stack trace
type GeneralError struct {
	Message   string
	Location  string
	Cause     error
	ErrorType int
}

func (err GeneralError) Error() string {
	if err.Cause == nil {
		return err.Message + " at " + err.Location
	}
	return err.Message + " at " + err.Location + "\n" + err.Cause.Error()
}

// CheckForMessage - Check if message is in stack trace
func (err *GeneralError) CheckForMessage(message string) bool {
	if strings.Contains(err.Error(), message) {
		return true
	}
	return false
}
