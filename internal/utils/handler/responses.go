package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.strale.io/go-travel/internal/database"
)

type ErrorResponseBody struct {
	Error string `json:"error"`
}

type Responder struct {
	origin string
}

func NewResponder(origin string) *Responder {
	return &Responder{
		origin: origin,
	}
}

func (r *Responder) ResolveErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", r.origin)
	switch {
	case errors.Is(err, database.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.As(err, &ErrBadRequest{}):
		w.WriteHeader(http.StatusBadRequest)
	case errors.As(err, &ErrForbidden{}):
		w.WriteHeader(http.StatusForbidden)
	case errors.As(err, &ErrUnauthorized{}):
		w.WriteHeader(http.StatusUnauthorized)
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

type Marshalable interface {
	MarshalJSON() ([]byte, error)
}

func (r *Responder) Respond(w http.ResponseWriter, status int, body Marshalable) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", r.origin)
	w.WriteHeader(status)
	if body != nil {
		bytesBody, err := body.MarshalJSON()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to serialize payload"}`))
			return
		}
		w.Write(bytesBody)
	}
}
