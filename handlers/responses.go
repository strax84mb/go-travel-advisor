package handlers

import (
	"net/http"
)

func writeBadRequest(w http.ResponseWriter, reason string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(reason))
}

func writeISE(w http.ResponseWriter, reason string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(reason))
}

func writeISEWithError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func writeISEWithReasonAndError(w http.ResponseWriter, reason string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(reason + "\n" + err.Error()))
}

func writeUnauthorised(w http.ResponseWriter, reason string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(reason))
}
