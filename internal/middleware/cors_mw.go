package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type corsMiddleware struct {
	origin string
}

func NewCorsMiddleware(origin string) *corsMiddleware {
	return &corsMiddleware{
		origin: origin,
	}
}

func (cmw *corsMiddleware) Middleware(httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Add("Access-Control-Allow-Origin", cmw.origin)
			w.Header().Add("Access-Control-Allow-Credentials", "false")
			w.Header().Add("Access-Control-Expose-Headers", "Content-Type,Origin,Accept")
			w.Header().Add("Access-Control-Max-Age", "900")
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Origin,Accept,Authorization,X-Requested-With")
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.Header().Add("Access-Control-Allow-Origin", cmw.origin)
		}
	})
}

func (cmw *corsMiddleware) AddOptionsHandlersForRoures(router *mux.Router) error {
	routesMap := make(map[string][]string)
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return err
		}
		if len(methods) > 0 {
			routesMap[path] = append(routesMap[path], methods...)
		}
		return nil
	})
	if err != nil {
		return err
	}
	var allMethods string
	for path, methods := range routesMap {
		allMethods = strings.Join(append(methods, http.MethodOptions), ",")
		router.Path(path).
			Methods(http.MethodOptions).
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Access-Control-Allow-Methods", allMethods)
				w.WriteHeader(http.StatusNoContent)
			})
		routesMap[path] = nil
	}
	return nil
}
