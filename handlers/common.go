package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func getBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading payload", http.StatusBadRequest)
		return false
	}
	if err = json.Unmarshal(bytes, v); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return false
	}
	return true
}

func checkJwt(w http.ResponseWriter, r *http.Request, role string) (string, bool) {
	token := r.Header.Get("Authorization")
	username, err := validateJwt(token, role)
	if err != nil {
		log.Printf("Error while authorizing: %s\n", err.Error())
		http.Error(w, "User is unauthorized", http.StatusUnauthorized)
		return "", false
	}
	return username, true
}

func serializeResponse(w http.ResponseWriter, payload interface{}) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error while marshaling response: %s\n", err.Error())
		http.Error(w, "Error while marshaling response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(payloadBytes)
}

func getInt64FromPath(w http.ResponseWriter, p httprouter.Params, varName string, errorText string) (int64, bool) {
	value, err := strconv.ParseInt(p.ByName(varName), 10, 64)
	if err != nil {
		http.Error(w, errorText, http.StatusBadRequest)
		return 0, false
	}
	return value, true
}

func getIntFromQuery(w http.ResponseWriter, r *http.Request, varName string, defaultValue int, errorText string) (int, bool) {
	valueString, ok := r.URL.Query()[varName]
	if ok {
		value, err := strconv.Atoi(valueString[0])
		if err != nil {
			http.Error(w, errorText, http.StatusBadRequest)
			return 0, false
		}
		return value, true
	}
	return defaultValue, true
}
