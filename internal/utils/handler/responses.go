package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/users"
)

type ErrorResponseBody struct {
	Error string `json:"error"`
}

func ResolveErrorResponse(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, database.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.As(err, &ErrBadRequest{}):
		w.WriteHeader(http.StatusBadRequest)
	case errors.As(err, &ErrForbidden{}):
		w.WriteHeader(http.StatusForbidden)
	case errors.As(err, &users.ErrUnauthorized{}):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	resp, err := json.Marshal(&ErrorResponseBody{
		Error: err.Error(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"failed to serialize error"}`))
		return
	}
	w.Write(resp)
}

func Respond(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	if body != nil {
		bytesBody, err := json.Marshal(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to serialize payload"}`))
			return
		}
		w.Write(bytesBody)
	}
}
