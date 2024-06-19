package handler

import (
	"fmt"
	"io"
	"net/http"
)

type Unmarshalable interface {
	UnmarshalJSON(input []byte) error
}

func GetBody(r *http.Request, target Unmarshalable) error {
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
