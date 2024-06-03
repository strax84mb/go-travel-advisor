package handler

import "fmt"

type AtomicInt64Validator func(name string, value int64) error

func IntMustBePositive(name string, value int64) error {
	if value <= 0 {
		return ErrBadRequest{
			message: fmt.Sprintf("%s must be positive", name),
		}
	}

	return nil
}

func IntMustBeZeroOrPositive(name string, value int64) error {
	if value < 0 {
		return ErrBadRequest{
			message: fmt.Sprintf("%s must be zero or positive", name),
		}
	}

	return nil
}

type AtomicStringValidator func(name string, value string) error

type AtomicValidator func(name string, value any) error

func IntIsPositive(name string, value any) error {
	if value.(int64) <= 0 {
		return ErrBadRequest{
			message: fmt.Sprintf("%s must be positive", name),
		}
	}

	return nil
}

func IntIsZeroOrPositive(name string, value any) error {
	if value.(int64) < 0 {
		return ErrBadRequest{
			message: fmt.Sprintf("%s must be zero or positive", name),
		}
	}

	return nil
}
