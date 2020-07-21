package handlers

import (
	"net/http"

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
	if /*Change to username*/ _, ok := checkJwt(w, r, db.UserRoleAdmin); !ok {
		return
	}
	payload := &SaveCommentPayload{}
	if !getBody(w, r, payload) {
		return
	}

}

// UpdateComment - updating coment
func UpdateComment(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

// DeleteComment - deleting coment
func DeleteComment(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}
