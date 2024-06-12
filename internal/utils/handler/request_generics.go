package handler

import (
	"fmt"
	"net/http"
	"strconv"
)

type AtomicValidatorGeneric[T Int64 | String] func(param string, val T) error

var GreaterThanZero AtomicValidatorGeneric[Int64] = func(param string, val Int64) error {
	value := int64(val)
	if value <= 0 {
		return ErrBadRequest{
			message: fmt.Sprintf("%s must be positive", param),
		}
	}
	return nil
}

type ParamValue[T Int64 | String] interface {
	FromString(val string) error
	*T
}

type Int64 int64

func (i64 *Int64) FromString(val string) error {
	i, err := strconv.ParseInt(val, 10, 64)
	if err == nil {
		*i64 = Int64(i)
		return nil
	} else {
		return err
	}
}

type String string

func (s *String) FromString(val string) error {
	*s = String(val)
	return nil
}

func QueryGen[
	T Int64 | String,
	PT ParamValue[T],
](
	r *http.Request,
	param string,
	mandatory bool,
	defaultValue T,
	validators ...AtomicValidatorGeneric[T],
) (T, error) {
	present := r.URL.Query().Has(param)
	strValue := r.URL.Query().Get(param)
	var (
		val T
		err error
	)
	if mandatory {
		if !present || strValue == "" {
			return val, mandatoryQueryNotPresent(param)
		}
	} else if !present {
		return defaultValue, nil
	}
	err = PT.FromString(&val, strValue)
	if err != nil {
		return val, ErrBadRequest{message: fmt.Sprintf("%s is malformed", param)}
	}

	for _, validator := range validators {
		if err = validator(param, val); err != nil {
			return val, err
		}
	}

	return val, nil
}
