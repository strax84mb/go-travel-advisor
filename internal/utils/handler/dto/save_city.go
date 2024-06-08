package dto

//go:generate ffjson -noencoder $GOFILE
type SaveCityDto struct {
	Name string `json:"name"`
}
