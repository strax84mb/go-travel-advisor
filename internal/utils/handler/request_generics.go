package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ParamValue[T Int64 | Int | String] interface {
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

type Int int

func (i *Int) FromString(val string) error {
	ival, err := strconv.Atoi(val)
	if err == nil {
		*i = Int(ival)
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

func Query[
	T Int64 | Int | String,
	PT ParamValue[T],
](
	r *http.Request,
	param string,
	mandatory bool,
	defaultValue T,
	validators ...AtomicValidator[T],
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

func Path[
	T Int | Int64 | String,
	PT ParamValue[T],
](
	r *http.Request,
	param string,
	validators ...AtomicValidator[T],
) (T, error) {
	var val T
	vars := mux.Vars(r)
	if vars == nil {
		return val, ErrBadRequest{message: "path parameters not found"}
	}
	strValue, found := vars[param]
	if !found {
		return val, ErrBadRequest{message: fmt.Sprintf("%s not present", param)}
	}
	err := PT.FromString(&val, strValue)
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
