package dto

import "gitlab.strale.io/go-travel/internal/database"

//go:generate ffjson -nodecoder $GOFILE
type AirportDtos struct {
	Items []AirportDto `json:"items"`
}

type AirportDto struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	CityID int64  `json:"cityId"`
}

func AirportToDto(airport database.Airport) *AirportDto {
	return &AirportDto{
		ID:     airport.ID,
		Name:   airport.Name,
		CityID: airport.CityID,
	}
}

func AirportsToDtos(airports []database.Airport) *AirportDtos {
	return &AirportDtos{
		Items: ConvertArray(airports, AirportToDto),
	}
}
