package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	db "gitlab.strale.io/go-travel/database"
	"gitlab.strale.io/go-travel/importing"
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
		handleErrors(w, err, err,
			errorHandling{
				err:    &db.ForbidenError{},
				status: http.StatusForbidden,
			})
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
		handleErrors(w, err, err,
			errorHandling{
				err:    &db.NotFoundError{},
				status: http.StatusNotFound,
			})
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
		handleErrors(w, err, err,
			errorHandling{
				err:    &db.NotFoundError{},
				status: http.StatusNotFound,
			})
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
	city, _, err := db.GetCityByID(id, maxComments)
	if err != nil {
		handleErrors(w, err, err,
			errorHandling{
				err:    &db.NotFoundError{},
				status: http.StatusNotFound,
			})
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	serializeResponse(w, cities)
}

const (
	cityName    = "name"
	cityCountry = "country"
)

// ImportCities - import cities from CSV
func ImportCities(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer r.Body.Close()
	if _, ok := checkJwt(w, r, db.UserRoleAdmin); !ok {
		return
	}
	var nameIndex, countryIndex int
	var ok bool
	if nameIndex, ok = getIntFromHeader(w, r, cityName, "Must have header \"name\" and it must be a nonnegative number!"); !ok {
		return
	}
	if countryIndex, ok = getIntFromHeader(w, r, cityCountry, "Must have header \"country\" and it must be a nonnegative number!"); !ok {
		return
	}
	mapping := importing.FieldMapping{
		cityName:    nameIndex,
		cityCountry: countryIndex,
	}
	im := importing.Importer{
		NumberOfChannels:  5,
		ChannelBufferSize: 50,
		Comma:             ',',
		Comment:           '"',
		Mapping:           mapping,
		Timestamp:         time.Now(),
		ParseRow:          parseCityRow,
		EntitySaver:       cityEntitySaver,
	}
	im.Parse(r.Body)
	w.WriteHeader(http.StatusNoContent)
	return
}

func parseCityRow(fields []string, mapping importing.FieldMapping) (interface{}, error) {
	city := db.City{
		Name:    fields[mapping[cityName]],
		Country: fields[mapping[cityCountry]],
	}
	defer func() {
		if err := recover(); err != nil {
			log.Println("[Parsing error]", err)
		}
	}()
	return city, nil
}

func cityEntitySaver(entity interface{}) error {
	city := entity.(db.City)
	return db.AddNewCity(city.Name, city.Country)
}
