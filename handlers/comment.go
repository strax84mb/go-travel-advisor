package handlers

import (
	"net/http"

	"gitlab.strale.io/go-travel/common"

	"github.com/julienschmidt/httprouter"
	db "gitlab.strale.io/go-travel/database"
)

// SaveCommentPayload - payload for saving or updating comment
type SaveCommentPayload struct {
	Text   string `json:"text"`
	CityID int64  `json:"city-id"`
}

// PostComment - post new coment
func PostComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	username, ok := checkJwt(w, r, db.UserRoleAny)
	if !ok {
		return
	}
	payload := &SaveCommentPayload{}
	if !getBody(w, r, payload) {
		return
	}
	err := db.AddComment(payload.Text, username, payload.CityID)
	if err != nil {
		if err.ErrorType > 400 && err.ErrorType < 500 {
			http.Error(w, err.Message, http.StatusNotFound)
			return
		}
		http.Error(w, err.Message, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdateComment - updating coment
func UpdateComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	username, ok := checkJwt(w, r, db.UserRoleAny)
	if !ok {
		return
	}
	payload := &SaveCommentPayload{}
	if !getBody(w, r, payload) {
		return
	}
	id, ok := getInt64FromPath(w, p, "id", "Illegal value for comment ID")
	if !ok {
		return
	}
	errRaw := db.UpdateComment(id, payload.Text, username, payload.CityID)
	if errRaw != nil {
		err := errRaw.(common.GeneralError)
		if err.ErrorType > 400 && err.ErrorType < 500 {
			http.Error(w, err.Message, http.StatusNotFound)
			return
		}
		http.Error(w, err.Message, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteComment - deleting coment
func DeleteComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	username, ok := checkJwt(w, r, db.UserRoleAny)
	if !ok {
		return
	}
	payload := &SaveCommentPayload{}
	if !getBody(w, r, payload) {
		return
	}
	id, ok := getInt64FromPath(w, p, "id", "Illegal value for comment ID")
	if !ok {
		return
	}
	err := db.DeleteComment(id, username)
	if err != nil {
		if err.ErrorType == common.UserNotAllowed {
			http.Error(w, err.Message, http.StatusUnauthorized)
			return
		}
		if err.ErrorType > 400 && err.ErrorType < 500 {
			http.Error(w, err.Message, http.StatusNotFound)
			return
		}
		http.Error(w, err.Message, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
