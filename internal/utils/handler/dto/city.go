package dto

import "gitlab.strale.io/go-travel/internal/database"

//go:generate ffjson -nodecoder $GOFILE
type CityDtos struct {
	Items []CityDto `json:"items"`
}

type CityDto struct {
	ID       int64            `json:"id"`
	Name     string           `json:"name"`
	Airports []CityAirportDto `json:"airports,omitempty"`
	//Comments []Comment `json:"foreignkey:CityID;association_foreignkey:ID"`
}

type CityAirportDto struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func CityToDto(city database.City) *CityDto {
	return &CityDto{
		ID:   city.ID,
		Name: city.Name,
		Airports: ConvertArray(
			city.Airports,
			func(airport database.Airport) *CityAirportDto {
				return &CityAirportDto{
					ID:   airport.ID,
					Name: airport.Name,
				}
			},
		),
	}
}

func CitiesToDtos(cities []database.City) *CityDtos {
	return &CityDtos{
		Items: ConvertArray(cities, CityToDto),
	}
}
