package handler

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func OptionsAllowedMethods(router *mux.Router, path string, forMethods ...string) {
	allMethods := strings.Join(append(forMethods, http.MethodOptions), ",")
	router.Path(path).Methods(http.MethodOptions).HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Methods", allMethods)
			w.WriteHeader(http.StatusNoContent)
		},
	)
}
