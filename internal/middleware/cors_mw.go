package middleware

import "net/http"

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
