package dto

//go:generate ffjson -noencoder $GOFILE
type SaveRouteDto struct {
	SourceID      int64 `json:"source"`
	DestinationID int64 `json:"destination"`
	Price         int32 `json:"price"`
}
