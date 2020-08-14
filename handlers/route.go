package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	db "gitlab.strale.io/go-travel/database"
	"gitlab.strale.io/go-travel/importing"
)

const (
	routeMappingSourceAirportID      = "source-airport-id"
	routeMappingDestinationAirportID = "destination-airport-id"
	routeMappingPrice                = "price"
)

// ImportRoutes - import routes data from CSV
func ImportRoutes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	if _, ok := checkJwt(w, r, db.UserRoleAdmin); !ok {
		return
	}
	var sourceIDIndex, destinationIDIndex, priceIndex int
	var ok bool
	if sourceIDIndex, ok = getIntFromHeader(w, r, routeMappingSourceAirportID,
		"Must have header \"source-airport-id\" and it must be a nonnegative number!"); !ok {
		return
	}
	if destinationIDIndex, ok = getIntFromHeader(w, r, routeMappingDestinationAirportID,
		"Must have header \"destination-airport-id\" and it must be a nonnegative number!"); !ok {
		return
	}
	if priceIndex, ok = getIntFromHeader(w, r, routeMappingPrice, "Must have header \"price\" and it must be a nonnegative number!"); !ok {
		return
	}
	mapping := importing.FieldMapping{
		routeMappingSourceAirportID:      sourceIDIndex,
		routeMappingDestinationAirportID: destinationIDIndex,
		routeMappingPrice:                priceIndex,
	}
	im := importing.Importer{
		NumberOfChannels:  5,
		ChannelBufferSize: 50,
		Comma:             ',',
		Comment:           '"',
		Mapping:           mapping,
		Timestamp:         time.Now(),
		ParseRow:          parseRouteRow,
		EntitySaver:       routeEntitySaver,
	}
	im.Parse(r.Body)
	w.WriteHeader(http.StatusNoContent)
	return
}

type routeImport struct {
	sourceID      int64
	destinationID int64
	price         float64
}

func parseRouteRow(fields []string, mapping importing.FieldMapping) (interface{}, error) {
	sourceID, err := strconv.ParseInt(fields[mapping[routeMappingSourceAirportID]], 10, 64)
	if err != nil {
		return routeImport{}, err
	}
	destinationID, err := strconv.ParseInt(fields[mapping[routeMappingDestinationAirportID]], 10, 64)
	if err != nil {
		return routeImport{}, err
	}
	price, err := strconv.ParseFloat(fields[mapping[routeMappingPrice]], 32)
	if err != nil {
		return routeImport{}, err
	}
	return routeImport{
		sourceID:      sourceID,
		destinationID: destinationID,
		price:         price,
	}, nil
}

func routeEntitySaver(entity interface{}) error {
	r := entity.(routeImport)
	return db.SaveRoute(r.sourceID, r.destinationID, r.price)
}

// PlanRoute - find cheapest route
func PlanRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	if _, ok := checkJwt(w, r, db.UserRoleAny); !ok {
		return
	}
	var start, end int64
	var ok bool
	if start, ok = getInt64FromQuery(w, r, "start", -1, "Bad value for start"); !ok {
		return
	}
	if end, ok = getInt64FromQuery(w, r, "end", -1, "Bad value for end"); !ok {
		return
	}
	pathDto, err := db.FindCheapesRoute(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	serializeResponse(w, pathDto)
}
