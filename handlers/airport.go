package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	db "gitlab.strale.io/go-travel/database"
	"gitlab.strale.io/go-travel/importing"
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

const (
	airportMappingAirportID   = "airport-id"
	airportMappingName        = "name"
	airportMappingCityName    = "city-name"
	airportMappingCityCountry = "city-country"
)

// ImportAirports - import airports from CSV
func ImportAirports(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	if _, ok := checkJwt(w, r, db.UserRoleAdmin); !ok {
		return
	}
	var airportIDIndex, nameIndex, cityNameIndex, cityCountryIndex int
	var ok bool
	if airportIDIndex, ok = getIntFromHeader(w, r, airportMappingAirportID, "Must have header \"airport-id\" and it must be a nonnegative number!"); !ok {
		return
	}
	if nameIndex, ok = getIntFromHeader(w, r, airportMappingName, "Must have header \"name\" and it must be a nonnegative number!"); !ok {
		return
	}
	if cityNameIndex, ok = getIntFromHeader(w, r, airportMappingCityName, "Must have header \"city-name\" and it must be a nonnegative number!"); !ok {
		return
	}
	if cityCountryIndex, ok = getIntFromHeader(w, r, airportMappingCityCountry, "Must have header \"city-country\" and it must be a nonnegative number!"); !ok {
		return
	}
	mapping := importing.FieldMapping{
		airportMappingAirportID:   airportIDIndex,
		airportMappingName:        nameIndex,
		airportMappingCityName:    cityNameIndex,
		airportMappingCityCountry: cityCountryIndex,
	}
	im := importing.Importer{
		NumberOfChannels:  5,
		ChannelBufferSize: 50,
		Comma:             ',',
		Comment:           '"',
		Mapping:           mapping,
		Timestamp:         time.Now(),
		ParseRow:          parseAirportRow,
		EntitySaver:       airportEntitySaver,
	}
	im.Parse(r.Body)
	w.WriteHeader(http.StatusNoContent)
	return
}

type airportImport struct {
	airportID   int64
	name        string
	cityName    string
	cityCountry string
}

func parseAirportRow(fields []string, mapping importing.FieldMapping) (interface{}, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[Parsing error]", err)
		}
	}()
	airportID, err := strconv.ParseInt(fields[mapping[airportMappingAirportID]], 10, 64)
	if err != nil {
		return airportImport{}, err
	}
	a := airportImport{
		airportID:   airportID,
		name:        fields[mapping[airportMappingName]],
		cityName:    fields[mapping[airportMappingCityName]],
		cityCountry: fields[mapping[airportMappingCityCountry]],
	}
	return a, nil
}

func airportEntitySaver(entity interface{}) error {
	a := entity.(airportImport)
	return db.ImportSingleAirport(a.airportID, a.name, a.cityName, a.cityCountry)
}
