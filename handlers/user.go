package handlers

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	db "gitlab.strale.io/go-travel/database"
)

// UsernamePassRequest - payload for signup and login
type UsernamePassRequest struct {
	Username string `validate:"nonzero,min=5,max=30,regexp=^[a-z]*$"`
	Password string `validate:"nonzero,min=5,max=30,regexp=^[a-zA-Z0-9]*$"`
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
		var e db.UsernameTakenError
		if errors.As(err, &e) {
			http.Error(w, e.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		var nfe *db.NotFoundError
		var ue *db.UnauthorizedError
		if errors.As(err, &nfe) {
			http.Error(w, nfe.Error(), http.StatusNotFound)
		} else if errors.As(err, &ue) {
			http.Error(w, ue.Error(), http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	// generate token
	token, err := generateJwt(payload.Username, user.Role, salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
