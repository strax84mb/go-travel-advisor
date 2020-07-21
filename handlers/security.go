package handlers

import (
	"strings"
	"time"

	"gitlab.strale.io/go-travel/common"
	"gitlab.strale.io/go-travel/database"

	jwt "github.com/dgrijalva/jwt-go"
)

func generateJwt(username string, role string, salt string) (string, error) {
	now := time.Now()
	exp := now.Add(3600 * 1000000000)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  username,
		"role": role,
		"nbf":  now.Unix(),
		"iat":  now.Unix(),
		"exp":  exp.Unix(),
	})
	tokenString, err := token.SignedString([]byte(salt))
	if err != nil {
		return "", &common.GeneralError{
			Cause:    err,
			Message:  "Failed to sign jwt token!",
			Location: "handlers.security.generateJwt",
		}
	}
	return tokenString, nil
}

func validRole(expected string, role string) bool {
	if expected == "any" {
		return true
	}
	return expected == role
}

// returns username
func validateJwt(header string, expectedRole string) (string, error) {
	if expectedRole == "" {
		return "", &common.GeneralError{
			Message:  "Required role is not declared!",
			Location: "handler.security.validateJwt",
		}
	}
	if !strings.HasPrefix(header, "Bearer ") {
		return "", &common.GeneralError{
			Message:  "Missing authentication token!",
			Location: "handler.security.validateJwt",
		}
	}
	token, err := jwt.Parse(strings.TrimPrefix(header, "Bearer "), func(token *jwt.Token) (interface{}, error) {
		claims := token.Claims.(jwt.MapClaims)
		username := claims["sub"].(string)
		salt, err := database.GetUserSaltByUsername(username)
		if err != nil {
			return []byte{}, err
		}
		return []byte(salt), nil
	})
	if err != nil {
		return "", &common.GeneralError{
			Message:  "Error while parsing JWT!",
			Location: "handlers.security.validateJwt",
			Cause:    err,
		}
	}
	claims := token.Claims.(jwt.MapClaims)
	username := claims["sub"].(string)
	role := claims["role"].(string)
	if !validRole(expectedRole, role) {
		return "", &common.GeneralError{
			Message:  "Incorrect role!",
			Location: "handler.security.validateJwt",
		}
	}
	return username, nil
}
