package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Value[V int | int64 | string] struct {
	val V
}

func (v *Value[V]) Val() V {
	return v.val
}

type PrimitiveValue[
	V int | int64 | string,
	T Value[V],
] interface {
	Val() V
	*T
}

func Int64(val string) (int64, error) {
	return strconv.ParseInt(val, 10, 64)
}

func Int(val string) (int, error) {
	return strconv.Atoi(val)
}

func String(val string) (string, error) {
	return val, nil
}

func Query[
	V int | int64 | string,
	T Value[V],
	PV PrimitiveValue[V, T],
](
	r *http.Request,
	parse func(string) (V, error),
	param string,
	mandatory bool,
	defaultValue V,
	validators ...AtomicValidator[V, T, PV],
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
		return T{val: defaultValue}, nil
	}
	v, err := parse(strValue)
	if err != nil {
		return val, ErrBadRequest{message: fmt.Sprintf("%s is malformed", param)}
	}
	val = T{val: v}

	for _, validator := range validators {
		if err = validator(param, val); err != nil {
			return val, err
		}
	}

	return val, nil
}

func Path[
	V int | int64 | string,
	T Value[V],
	PV PrimitiveValue[V, T],
](
	r *http.Request,
	parse func(string) (V, error),
	param string,
	validators ...AtomicValidator[V, T, PV],
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
	v, err := parse(strValue)
	if err != nil {
		return val, ErrBadRequest{message: fmt.Sprintf("%s is malformed", param)}
	}
	val = T{val: v}

	for _, validator := range validators {
		if err = validator(param, val); err != nil {
			return val, err
		}
	}
	return val, nil
}
