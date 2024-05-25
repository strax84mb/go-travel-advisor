package middleware

import (
	"net/http"

	"gitlab.strale.io/go-travel/internal/utils"
	"gitlab.strale.io/go-travel/internal/utils/handler"
)

type securityService interface {
	VerifyJWT(r *http.Request) (int64, []string, error)
}

type verifyJWTMiddleware struct {
	securitySrvc securityService
}

func NewVerifyJWTMiddleware(securitySrvc securityService) *verifyJWTMiddleware {
	return &verifyJWTMiddleware{
		securitySrvc: securitySrvc,
	}
}

func (mw *verifyJWTMiddleware) Middleware(httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, roles, err := mw.securitySrvc.VerifyJWT(r)
		if err != nil {
			handler.ResolveErrorResponse(w, err)
			return
		}
		if id != 0 {
			ctx := r.Context()
			ctx = utils.WithValue(ctx, "userId", id)
			ctx = utils.WithValue(ctx, "userRoles", roles)
			r = r.Clone(ctx)
		}
		httpHandler.ServeHTTP(w, r)
	})
}
