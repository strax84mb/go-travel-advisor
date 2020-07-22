package handlers

import (
	"net/http"

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
		if err.ErrorType > 300 && err.ErrorType < 400 {
			http.Error(w, err.Message, http.StatusUnauthorized)
		} else {
			http.Error(w, err.Message, http.StatusInternalServerError)
		}
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
	if err := db.SaveNewUser(payload.Username, payload.Password); err != nil {
		if err.ErrorType > 300 && err.ErrorType < 400 {
			http.Error(w, err.Message, http.StatusUnauthorized)
		} else {
			http.Error(w, err.Message, http.StatusInternalServerError)
		}
		return
	}
	// generate token
	token, err := generateJwt(payload.Username, user.Role, salt)
	if err != nil {
		http.Error(w, err.Message, http.StatusInternalServerError)
		return
	}
	// respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
