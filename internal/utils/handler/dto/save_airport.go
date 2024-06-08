package dto

//go:generate ffjson -noencoder $GOFILE
type SaveAirportDto struct {
	Name   string `json:"name"`
	CityID int64  `json:"cityId"`
}
