package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func QueryAsInt(r *http.Request, param string, mandatory bool, defaultValue int, validators ...AtomicValidator) (int, error) {
	present := r.URL.Query().Has(param)
	strValue := r.URL.Query().Get(param)
	if mandatory {
		if !present {
			return 0, mandatoryQueryNotPresent(param)
		}
	} else if !present {
		return defaultValue, nil
	}
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, ErrBadRequest{message: fmt.Sprintf("%s is malformed", param)}
	}
	for _, validator := range validators {
		if err = validator(param, int64(intValue)); err != nil {
			return 0, err
		}
	}
	return intValue, nil
}

func PathAsInt64(r *http.Request, param string, validators ...AtomicValidator) (int64, error) {
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

func GetBody(r *http.Request, target interface{}) error {
	body := r.Body
	defer body.Close()
	bytes, err := io.ReadAll(body)
	if err != nil {
		return ErrBadRequest{message: fmt.Sprintf("payload cannot be extracted: %s", err.Error())}
	}
	if err = json.Unmarshal(bytes, target); err != nil {
		return ErrBadRequest{message: fmt.Sprintf("payload is malformed: %s", err.Error())}
	}
	return nil
}
