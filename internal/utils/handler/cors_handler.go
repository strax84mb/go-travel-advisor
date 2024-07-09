package handler

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Cors struct {
	origin string
}

func NewCors(origin string) *Cors {
	return &Cors{
		origin: origin,
	}
}

func (c *Cors) Options(router *mux.Router, path string, methods ...string) {
	router.Path(path).Methods(methods...).HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", c.origin)
			w.Header().Add("Access-Control-Allow-Credentials", "true")
			w.Header().Add("Access-Control-Expose-Headers", "Content-Type, Origin, Accept")
			w.Header().Add("Access-Control-Max-Age", "900")
			w.Header().Add("Access-Control-Allow-Methods", strings.Join(methods, ","))
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, Authorization")
			w.WriteHeader(http.StatusOK)
		},
	)
}
