package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	db "gitlab.strale.io/go-travel/database"
)

// ListAirports - list all airports
func ListAirports(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	if _, ok := checkJwt(w, r, db.UserRoleAny); !ok {
		return
	}
	airports, err := db.ListAirports()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	serializeResponse(w, airports)
}

// WriteAirportPayload - payload for adding and updating airport
type WriteAirportPayload struct {
	AirportID int64
	Name      string
	CityID    int64
}

// AddAirport - save new airport
func AddAirport(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	if _, ok := checkJwt(w, r, db.UserRoleAdmin); !ok {
		return
	}
	payload := WriteAirportPayload{}
	if !getBody(w, r, &payload) {
		return
	}
	airport, err := db.SaveAirport(payload.AirportID, payload.Name, payload.CityID)
	if err != nil {
		handleErrors(w, err, err,
			errorHandling{
				err:    &db.ForbidenError{},
				status: http.StatusForbidden,
			})
		return
	}
	serializeResponse(w, airport)
}

// GetAirport - get airport with specific ID
func GetAirport(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer r.Body.Close()
	if _, ok := checkJwt(w, r, db.UserRoleAdmin); !ok {
		return
	}
	var id int64
	var ok bool
	if id, ok = getInt64FromPath(w, p, "id", "Bad value of ID"); !ok {
		return
	}
	airport, err := db.GetAirport(id)
	if err != nil {
		handleErrors(w, err, err,
			errorHandling{
				err:    &db.NotFoundError{},
				status: http.StatusNotFound,
			})
		return
	}
	serializeResponse(w, airport)
}

// UpdateAirport - change data of a specific airport
func UpdateAirport(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer r.Body.Close()
	if _, ok := checkJwt(w, r, db.UserRoleAdmin); !ok {
		return
	}
	var id int64
	var ok bool
	if id, ok = getInt64FromPath(w, p, "id", "Bad value of ID"); !ok {
		return
	}
	payload := WriteAirportPayload{}
	if !getBody(w, r, &payload) {
		return
	}
	airport, err := db.UpdateAirport(id, payload.AirportID, payload.Name, payload.CityID)
	if err != nil {
		handleErrors(w, err, err,
			errorHandling{
				err:    &db.NotFoundError{},
				status: http.StatusNotFound,
			},
			errorHandling{
				err:    &db.ForbidenError{},
				status: http.StatusForbidden,
			})
		return
	}
	serializeResponse(w, airport)
}

// DeleteAirport - delete and airport with given ID
func DeleteAirport(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer r.Body.Close()
	if _, ok := checkJwt(w, r, db.UserRoleAdmin); !ok {
		return
	}
	var id int64
	var ok bool
	if id, ok = getInt64FromPath(w, p, "id", "Bad value of ID"); !ok {
		return
	}
	err := db.DeleteAirport(id)
	if err != nil {
		handleErrors(w, err, err,
			errorHandling{
				err:    &db.NotFoundError{},
				status: http.StatusNotFound,
			})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ImportAirports - import airports from CSV
func ImportAirports(w http.ResponseWriter, r *http.Request, p httprouter.Params) {}
