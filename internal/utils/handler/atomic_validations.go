package handler

import "fmt"

type AtomicValidator[
	V int | int64 | string,
	T Value[V],
	PV PrimitiveValue[V, T],
] func(param string, val T) error

//type NumValidator[T Int64 | Int] AtomicValidator[T]

func IsPositive[
	V int | int64,
	T Value[V],
	PV PrimitiveValue[V, T],
](param string, val T) error {
	v := PV.Val(&val)
	if v <= 0 {
		return ErrBadRequest{
			message: fmt.Sprintf("%s must be positive", param),
		}
	}
	return nil
}

func IsZeroOrPositive[
	V int | int64,
	T Value[V],
	PV PrimitiveValue[V, T],
](param string, val T) error {
	v := PV.Val(&val)
	if v < 0 {
		return ErrBadRequest{
			message: fmt.Sprintf("%s must be zero or positive", param),
		}
	}
	return nil
}

//type AtomicStringValidator AtomicValidator[String]

func IsNotEmpty(param string, value string) error {
	if value == "" {
		return ErrBadRequest{
			message: param + " must not be empty",
		}
	}
	return nil
}
