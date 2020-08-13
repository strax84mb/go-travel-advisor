package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/pquerna/ffjson/ffjson"
	v2 "gopkg.in/validator.v2"
)

func getBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading payload", http.StatusBadRequest)
		return false
	}
	if err = ffjson.Unmarshal(bytes, v); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return false
	}
	if err = v2.Validate(v); err != nil {
		http.Error(w, "Invalid payload"+err.Error(), http.StatusBadRequest)
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

func getInt64FromQuery(w http.ResponseWriter, r *http.Request, varName string, defaultValue int64, errorText string) (int64, bool) {
	valueString, ok := r.URL.Query()[varName]
	if ok {
		value, err := strconv.ParseInt(valueString[0], 10, 64)
		if err != nil {
			http.Error(w, errorText, http.StatusBadRequest)
			return 0, false
		}
		return value, true
	}
	return defaultValue, true
}

func getIntFromHeader(w http.ResponseWriter, r *http.Request, varName string, errorText string) (int, bool) {
	valueString := r.Header.Get(varName)
	if valueString == "" {
		http.Error(w, errorText, http.StatusBadRequest)
		return 0, false
	}
	value, err := strconv.Atoi(valueString)
	if err != nil || value < 0 {
		log.Printf("Error while converting header value %s to int! Error: %s", varName, err.Error())
		http.Error(w, errorText, http.StatusBadRequest)
		return 0, false
	}
	return value, true
}

type errorHandling struct {
	err               error
	status            int
	substituteMessage string
}

func handleErrors(w http.ResponseWriter, defaultError error, chainError error, handlings ...errorHandling) {
	if handlings != nil || len(handlings) > 0 {
		for _, h := range handlings {
			if errors.As(chainError, &h.err) {
				if h.substituteMessage == "" {
					http.Error(w, h.err.Error(), h.status)
				} else {
					http.Error(w, h.substituteMessage, h.status)
				}
				return
			}
		}
	}
	http.Error(w, defaultError.Error(), http.StatusInternalServerError)
}
