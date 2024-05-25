package dto

import "gitlab.strale.io/go-travel/internal/database"

type AirportDto struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	CityID int64  `json:"cityId"`
}

func AirportToDto(airport database.Airport) AirportDto {
	return AirportDto{
		ID:     airport.ID,
		Name:   airport.Name,
		CityID: airport.CityID,
	}
}

type SaveAirportDto struct {
	Name   string `json:"name"`
	CityID int64  `json:"cityId"`
}
