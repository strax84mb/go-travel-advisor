package dto

import "gitlab.strale.io/go-travel/internal/database"

type RouteAirportDto struct {
	ID   int64  `json:"id"`
	Name string `json:"name,omitempty"`
}

type RouteDto struct {
	ID          int64           `json:"id"`
	Source      RouteAirportDto `json:"source"`
	Destination RouteAirportDto `json:"destination"`
	Price       float32         `json:"price"`
}

func RouteToDto(route database.Route) RouteDto {
	source := RouteAirportDto{
		ID: route.SourceID,
	}
	if route.Source != nil {
		source.Name = route.Source.Name
	}
	destination := RouteAirportDto{
		ID: route.DestinationID,
	}
	if route.Destination != nil {
		destination.Name = route.Destination.Name
	}
	return RouteDto{
		ID:          route.ID,
		Source:      source,
		Destination: destination,
		Price:       route.Price,
	}
}

func RoutesToDtos(routes []database.Route) []RouteDto {
	list := make([]RouteDto, len(routes))
	for i, route := range routes {
		list[i] = RouteToDto(route)
	}
	return list
}

type SaveRouteDto struct {
	SourceID      int64   `json:"source"`
	DestinationID int64   `json:"destination"`
	Price         float32 `json:"price"`
}

type UpdateRoutePrice struct {
	Price float32 `json:"price"`
}
