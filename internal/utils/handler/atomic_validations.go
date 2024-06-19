package handler

import "fmt"

type AtomicValidator[T Int64 | Int | String] func(param string, val T) error

type NumValidator[T Int64 | Int] AtomicValidator[T]

func IsPositive[T Int64 | Int](param string, val T) error {
	if val <= 0 {
		return ErrBadRequest{
			message: fmt.Sprintf("%s must be positive", param),
		}
	}
	return nil
}

func IsZeroOrPositive[T Int64 | Int](param string, val T) error {
	if val < 0 {
		return ErrBadRequest{
			message: fmt.Sprintf("%s must be zero or positive", param),
		}
	}
	return nil
}

type AtomicStringValidator AtomicValidator[String]

func IsNotEmpty(param string, value string) error {
	if value == "" {
		return ErrBadRequest{
			message: param + " must not be empty",
		}
	}
	return nil
}
