package dto

//go:generate ffjson -noencoder $GOFILE
type UpdateRoutePrice struct {
	Price float32 `json:"price"`
}
