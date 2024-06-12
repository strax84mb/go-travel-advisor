package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func QueryAsInt(r *http.Request, param string, mandatory bool, defaultValue int, validators ...AtomicInt64Validator) (int, error) {
	value, err := query[int](r, param, true, defaultValue)
	if err != nil {
		return 0, err
	}
	for _, validator := range validators {
		if err = validator(param, int64(value)); err != nil {
			return 0, err
		}
	}
	return value, nil
}

func QueryAsInt64(r *http.Request, param string, mandatory bool, defaultValue int64, validators ...AtomicInt64Validator) (int64, error) {
	value, err := query[int64](r, param, true, defaultValue)
	if err != nil {
		return 0, err
	}
	for _, validator := range validators {
		if err = validator(param, value); err != nil {
			return 0, err
		}
	}
	return value, nil
}

func QueryAsString(r *http.Request, param string, mandatory bool, defaultValue string, validators ...AtomicStringValidator) (string, error) {
	value, err := query[string](r, param, true, defaultValue)
	if err != nil {
		return "", err
	}
	for _, validator := range validators {
		if err := validator(param, value); err != nil {
			return "", err
		}
	}
	return value, nil
}

func query[T int | int64 | string](
	r *http.Request,
	param string,
	mandatory bool,
	defaultValue T,
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
	switch t := any(val).(type) {
	case string:
		val = any(strValue).(T)
	case int:
		t, err = strconv.Atoi(strValue)
		if err != nil {
			return val, ErrBadRequest{message: fmt.Sprintf("%s is malformed", param)}
		}
		val = any(t).(T)
	case int64:
		t, err = strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			return val, ErrBadRequest{message: fmt.Sprintf("%s is malformed", param)}
		}
		val = any(t).(T)
	}
	return val, nil
}

func PathAsInt64(r *http.Request, param string, validators ...AtomicInt64Validator) (int64, error) {
	vars := mux.Vars(r)
	if vars == nil {
		return 0, ErrBadRequest{message: "path parameters not found"}
	}
	strValue, found := vars[param]
	if !found {
		return 0, ErrBadRequest{message: fmt.Sprintf("%s not present", param)}
	}
	intValue, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		return 0, ErrBadRequest{message: fmt.Sprintf("parameter %s is malformed", param)}
	}
	for _, validator := range validators {
		if err = validator(param, intValue); err != nil {
			return 0, err
		}
	}
	return intValue, nil
}

type Unmarshalable interface {
	UnmarshalJSON(input []byte) error
}

func GetBodyFF(r *http.Request, target Unmarshalable) error {
	body := r.Body
	defer body.Close()
	bytes, err := io.ReadAll(body)
	if err != nil {
		return ErrBadRequest{message: fmt.Sprintf("payload cannot be extracted: %s", err.Error())}
	}
	if err = target.UnmarshalJSON(bytes); err != nil {
		return ErrBadRequest{message: fmt.Sprintf("payload is malformed: %s", err.Error())}
	}
	return nil
}
