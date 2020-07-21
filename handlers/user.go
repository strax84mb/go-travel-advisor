package handlers

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	db "gitlab.strale.io/go-travel/database"
)

// UsernamePassRequest - payload for signup and login
type UsernamePassRequest struct {
	Username string
	Password string
}

// SignupUser - handler for endpoint /user/signup
func SignupUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	payload := &UsernamePassRequest{}
	if !getBody(w, r, payload) {
		return
	}
	// save new user
	if err := db.SaveNewUser(payload.Username, payload.Password); err != nil {
		writeISEWithError(w, err)
		return
	}
	// respond with success
	w.WriteHeader(http.StatusNoContent)
}

// LoginUser - handler for endpoint /user/login
func LoginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	payload := &UsernamePassRequest{}
	if !getBody(w, r, payload) {
		return
	}
	// validate user
	user, salt, err := db.GetUserByUsernameAndPassword(payload.Username, payload.Password)
	if err != nil {
		if strings.Contains(err.Error(), "Incorrect password!") || strings.Contains(err.Error(), "Incorrect username!") {
			writeUnauthorised(w, "Username or password are incorrect!")
		} else {
			writeISEWithError(w, err)
		}
		return
	}
	// generate token
	token, err := generateJwt(payload.Username, user.Role, salt)
	if err != nil {
		writeISE(w, "Error while generating JWT!")
		return
	}
	// respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
