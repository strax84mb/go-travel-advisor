package database

// NotFoundError - when record is not found
type NotFoundError struct {
	Message string
}

func (err NotFoundError) Error() string {
	return err.Message
}

// StatementError - when error happened statement execution
type StatementError struct {
	Message string
	Cause   error
}

func (err StatementError) Error() string {
	return err.Message
}

func (err StatementError) Unwrap() error {
	return err.Cause
}

// UnauthorizedError - error when unauthorized action was attempted
type UnauthorizedError struct {
	Message string
	Cause   error
}

func (e UnauthorizedError) Error() string {
	return e.Error()
}

// UsernameTakenError - when username is already taken
type UsernameTakenError struct{}

func (u UsernameTakenError) Error() string {
	return "Username is taken!"
}

// ForbidenError - when action is not allowed for any reason
type ForbidenError struct {
	Message string
}

func (e ForbidenError) Error() string {
	return e.Message
}
