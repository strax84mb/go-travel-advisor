package handlers

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	db "gitlab.strale.io/go-travel/database"
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
		log.Printf("Error while generating JWT! Error: %s\n", err.Error())
		return "", errors.New("Failed to sign jwt token")
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
		return "", errors.New("Required role is not declared")
	}
	if !strings.HasPrefix(header, "Bearer ") {
		return "", errors.New("Missing authentication token")
	}
	token, err := jwt.Parse(strings.TrimPrefix(header, "Bearer "), func(token *jwt.Token) (interface{}, error) {
		claims := token.Claims.(jwt.MapClaims)
		username := claims["sub"].(string)
		salt, err := db.GetUserSaltByUsername(username)
		if err != nil {
			return []byte{}, err
		}
		return []byte(salt), nil
	})
	if err != nil {
		return "", errors.New("Error while parsing JWT")
	}
	claims := token.Claims.(jwt.MapClaims)
	username := claims["sub"].(string)
	role := claims["role"].(string)
	if !validRole(expectedRole, role) {
		return "", errors.New("Incorrect role")
	}
	return username, nil
}
