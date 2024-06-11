package dto

import (
	"bytes"
	"fmt"

	"github.com/pquerna/ffjson/ffjson"
	"gitlab.strale.io/go-travel/internal/database"
)

type CheapestPath struct {
	Path      []Step `json:"path"`
	FullPrice int64  `json:"fullPrice"`
}

func routeToFlight(route *database.Route) *FlightDto {
	return &FlightDto{
		Price: route.Price,
		From: FlightAirportDto{
			ID:   route.SourceID,
			Name: route.Source.Name,
			City: FlightCityDto{
				ID:   route.Source.CityID,
				Name: route.Source.City.Name,
			},
		},
		To: FlightAirportDto{
			ID:   route.DestinationID,
			Name: route.Destination.Name,
			City: FlightCityDto{
				ID:   route.Destination.CityID,
				Name: route.Destination.City.Name,
			},
		},
	}
}

func constructTransfer(prevRoute *database.Route, nextRoute *database.Route) *TransferDto {
	return &TransferDto{
		City: TransferCityDto{
			ID:   nextRoute.Source.CityID,
			Name: nextRoute.Source.City.Name,
		},
		From: TransferAirportDto{
			ID:   prevRoute.DestinationID,
			Name: prevRoute.Destination.Name,
		},
		To: TransferAirportDto{
			ID:   nextRoute.SourceID,
			Name: nextRoute.Source.Name,
		},
	}
}

func CompileCheapestPath(routes []*database.Route, fullPrice int64) *CheapestPath {
	if routes == nil {
		return nil
	}
	var path []Step
	path = append(path, routeToFlight(routes[0]))
	if len(routes) == 1 {
		return &CheapestPath{
			Path:      path,
			FullPrice: fullPrice,
		}
	}
	var prevRoute *database.Route
	for i, route := range routes {
		if i == 0 {
			continue
		}
		prevRoute = routes[i]
		if prevRoute.DestinationID != route.SourceID {
			path = append(path, constructTransfer(prevRoute, route))
		}
		path = append(path, routeToFlight(route))
	}
	return &CheapestPath{
		Path:      path,
		FullPrice: fullPrice,
	}
}

func (cp *CheapestPath) Encode() ([]byte, error) {
	var err error
	buf := new(bytes.Buffer)
	enc := ffjson.NewEncoder(buf)
	if _, err = buf.Write([]byte(`{"path":[`)); err != nil {
		return nil, err
	}
	cap := len(cp.Path) - 1
	for i, step := range cp.Path {
		if err = enc.Encode(step); err != nil {
			return nil, err
		}
		if i < cap {
			if _, err = buf.Write([]byte(",")); err != nil {
				return nil, err
			}
		}
	}
	//TODO
	_, err = buf.WriteString(fmt.Sprintf(`],"fullPrice":%d}`, cp.FullPrice))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
