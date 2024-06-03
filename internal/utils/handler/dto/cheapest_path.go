package dto

import "gitlab.strale.io/go-travel/internal/database"

type StepType string

const (
	FlightType   StepType = "FLIGHT"
	TransferType StepType = "TRANSFER"
)

type FlightCityDto struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type FlightAirportDto struct {
	ID   int64         `json:"id"`
	Name string        `json:"name"`
	City FlightCityDto `json:"city"`
}

type FlightDto struct {
	From  FlightAirportDto `json:"from"`
	To    FlightAirportDto `json:"to"`
	Price float32          `json:"price"`
}

func (f *FlightDto) Type() StepType {
	return FlightType
}

type TransferAirportDto struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TransferCityDto struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TransferDto struct {
	City TransferCityDto    `json:"city"`
	From TransferAirportDto `json:"from"`
	To   TransferAirportDto `json:"to"`
}

func (t *TransferDto) Type() StepType {
	return TransferType
}

type Step interface {
	Type() StepType
}

type CheapestPath struct {
	Path      []Step  `json:"path"`
	FullPrice float32 `json:"fullPrice"`
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

func CompileCheapestPath(routes []*database.Route, fullPrice float32) *CheapestPath {
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
