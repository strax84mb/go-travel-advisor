package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	db "gitlab.strale.io/go-travel/database"
)

// WriteCityPayload - payload for adding and updating a city
type WriteCityPayload struct {
	Name    string
	Country string
}

// AddCity - add new city
func AddCity(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	if _, ok := checkJwt(w, r, db.UserRoleAdmin); !ok {
		return
	}
	payload := &WriteCityPayload{}
	if !getBody(w, r, payload) {
		return
	}
	err := db.AddNewCity(payload.Name, payload.Country)
	if err != nil {
		log.Printf("Error while adding new city: %s", err.Error())
		http.Error(w, "Error while adding new city", http.StatusInternalServerError)
		return
	}
	// respond with success
	w.WriteHeader(http.StatusNoContent)
}

// UpdateCity - update city
func UpdateCity(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer r.Body.Close()
	if _, ok := checkJwt(w, r, db.UserRoleAdmin); !ok {
		return
	}
	var id int64
	var ok bool
	if id, ok = getInt64FromPath(w, p, "id", "Bad value of ID"); !ok {
		return
	}
	payload := &WriteCityPayload{}
	if !getBody(w, r, payload) {
		return
	}
	err := db.UpdateCity(id, payload.Name, payload.Country)
	if err != nil {
		log.Printf("Error while updating city: %s\n", err.Error())
		http.Error(w, "Error while updating city", http.StatusInternalServerError)
		return
	}
	// respond with success
	w.WriteHeader(http.StatusNoContent)
}

// DeleteCity - delete city
func DeleteCity(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer r.Body.Close()
	if _, ok := checkJwt(w, r, db.UserRoleAdmin); !ok {
		return
	}
	var id int64
	var ok bool
	if id, ok = getInt64FromPath(w, p, "id", "Bad value of ID"); !ok {
		return
	}
	err := db.DeleteCity(id)
	if err != nil {
		writeISEWithReasonAndError(w, "Error while deleting city", err)
		return
	}
	// respond with success
	w.WriteHeader(http.StatusNoContent)
}

// GetCity - get specific city
func GetCity(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer r.Body.Close()
	if _, ok := checkJwt(w, r, db.UserRoleAny); !ok {
		return
	}
	var id int64
	var ok bool
	if id, ok = getInt64FromPath(w, p, "id", "Bad value of ID"); !ok {
		return
	}
	var maxComments int
	if maxComments, ok = getIntFromQuery(w, r, "max-comments", -1, "Bad value for max-comments"); !ok {
		return
	}
	city, found, err := db.GetCityByID(id, maxComments)
	if err != nil {
		log.Printf("Error while reading city: %s\n", err.Error())
		http.Error(w, "Error while reading city", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, fmt.Sprintf("City with ID %d not found!", id), http.StatusNotFound)
		return
	}
	serializeResponse(w, city)
}

// GetAllCities - list all cities
func GetAllCities(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	if _, ok := checkJwt(w, r, db.UserRoleAny); !ok {
		return
	}
	var maxComments int
	var ok bool
	if maxComments, ok = getIntFromQuery(w, r, "max-comments", -1, "Bad value for max-comments"); !ok {
		return
	}
	cities, err := db.GetAllCities(maxComments)
	if err != nil {
		log.Printf("Error while listing all cities: %s\n", err.Error())
		http.Error(w, "Error while listing all cities", http.StatusInternalServerError)
		return
	}
	serializeResponse(w, cities)
}

// ImportCities - import cities from CSV
func ImportCities(w http.ResponseWriter, r *http.Request, p httprouter.Params) {}
