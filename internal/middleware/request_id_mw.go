package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"gitlab.strale.io/go-travel/internal/utils"
)

func RequestIDMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := uuid.New().String()
		ctx := r.Context()
		ctx = utils.WithValue(ctx, "requestId", key)
		r = r.Clone(ctx)
		handler.ServeHTTP(w, r)
	})
}
