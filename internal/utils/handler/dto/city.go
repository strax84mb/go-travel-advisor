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
	var airports []CityAirportDto
	if city.Airports != nil {
		airports = make([]CityAirportDto, len(city.Airports))
		for i, airport := range city.Airports {
			airports[i] = CityAirportDto{
				ID:   airport.ID,
				Name: airport.Name,
			}
		}
	}
	return &CityDto{
		ID:       city.ID,
		Name:     city.Name,
		Airports: airports,
	}
}

func CitiesToDtos(cities []database.City) *CityDtos {
	dtos := CityDtos{
		Items: make([]CityDto, len(cities)),
	}
	for i, c := range cities {
		dtos.Items[i] = *CityToDto(c)
	}
	return &dtos
}
