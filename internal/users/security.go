package users

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/utils/handler"
	"golang.org/x/crypto/ssh"
)

type securityService struct {
	key interface{}
}

func NewSecurityService(key string) (*securityService, error) {
	rsaKey, err := ssh.ParseRawPrivateKey([]byte(key))
	if err != nil {
		return nil, err
	}
	return &securityService{
		key: rsaKey,
	}, nil
}

func (ss *securityService) GenerateJWT(ctx context.Context, user database.User) (string, error) {
	rolesStr := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		rolesStr[i] = string(role.Role)
	}
	now := time.Now()
	exp := now.Add(3600 * 1000000000)
	token := jwt.NewWithClaims(
		jwt.SigningMethodPS256.SigningMethodRSA,
		jwt.MapClaims{
			"sub":      user.Username,
			"username": user.Username,
			"user-id":  user.ID,
			"roles":    strings.Join(rolesStr, ","),
			"nbf":      now.Unix(),
			"iat":      now.Unix(),
			"exp":      exp.Unix(),
		},
	)
	tokenString, err := token.SignedString(ss.key)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("error creating jwt")
		return "", fmt.Errorf("error creating jwt: %w", err)
	}
	return tokenString, nil
}

func (ss *securityService) VerifyJWT(r *http.Request) (int64, []string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// assumed to be accessible to all
		return 0, nil, nil
	}
	if strings.HasPrefix(authHeader, "Bearer ") {
		return 0, nil, handler.NewErrUnauthorizedWithCause("malformed authorization header")
	}
	authHeader = strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(
		authHeader,
		func(t *jwt.Token) (interface{}, error) {
			return ss.key, nil
		},
	)
	if err != nil {
		return 0, nil, handler.NewErrUnauthorizedWithCause(err.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	roles := strings.Split(claims["roles"].(string), ",")
	return claims["user-id"].(int64), roles, nil
}
