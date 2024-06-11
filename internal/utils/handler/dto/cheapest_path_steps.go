package dto

//go:generate ffjson -nodecoder $GOFILE

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
	Price int32            `json:"price"`
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
